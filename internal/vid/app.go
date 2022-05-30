package vid

import (
	"context"
	"fmt"
	"github.com/torbenschinke/picd/internal/pic"
	"github.com/torbenschinke/picd/internal/pic/raspistill"
	"github.com/torbenschinke/picd/internal/sensor"
	"github.com/torbenschinke/picd/pkg/server"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Application struct {
	dir          string
	timelapseDir string
	ticker       *time.Ticker
}

func NewApplication(ctx context.Context) *Application {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir := filepath.Join(home, "picd")
	timelapseDir := filepath.Join(home, "picd", "timelapse")

	_ = os.MkdirAll(timelapseDir, os.ModePerm)

	return &Application{
		dir:          dir,
		timelapseDir: timelapseDir,
		ticker:       time.NewTicker(time.Hour),
	}
}

func (a *Application) Run(ctx context.Context) error {
	camRepo := raspistill.NewCameraRepo()
	picService := pic.NewCameraService(camRepo)
	sensorService := sensor.NewSenseService()

	go a.runTimelapse()
	go a.removeOldest()

	return server.Serve(ctx, NewController(a, sensorService, picService))
}

func (a *Application) removeOldest() {
	for range a.ticker.C {
		files, err := os.ReadDir(a.timelapseDir)
		if err != nil {
			panic(err)
		}

		now := time.Now()
		count := 0
		for _, file := range files {
			fname := filepath.Join(a.timelapseDir, file.Name())
			info, err := file.Info()
			if err != nil {
				fmt.Printf("cannot get fileinfo: %s: %v\n", fname, info)
			}

			if now.Sub(info.ModTime()).Hours() > 24 {
				if err := os.Remove(fname); err != nil {
					fmt.Printf("failed to remove file older than 24h %s: %v\n", fname, err)
				}
				count++
			}
		}

		fmt.Printf("purged %d files, still having %d\n", count, len(files)-count)
	}
}

func (a *Application) runTimelapse() {
	for {
		// https://www.waveshare.com/wiki/RPi_Camera_(C)#libcamera-still
		// libcamera-still --timelapse 1000 -q 80 --timestamp  --latest current.jpg --nopreview --width 1280 --height 720 -t 0
		cmd := exec.Command("/usr/bin/libcamera-still",
			"--timelapse", "1000",
			"-q", "80",
			"--timestamp",
			"--latest", "current.jpg",
			"--nopreview",
			"--width", "960", // different aspect ratio => use entire sensor
			"--height", "720",
			"-t", "0",
			"--exif", "IFD0.Orientation=8", // 3=180, 6=90, 8=270
			"--saturation", "0",
		)

		cmd.Env = os.Environ()
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard //os.Stdout // very chatty, cannot turn that off
		cmd.Dir = a.timelapseDir

		if err := cmd.Run(); err != nil {
			fmt.Printf("libcamera-still timelapse failed: %v\n", err)
			time.Sleep(5 * time.Second)
			fmt.Printf("retry libcamera\n")
			continue
		}

		break // exit without error, sighub?
	}
}
