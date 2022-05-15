package raspistill

import (
	"bytes"
	"fmt"
	"github.com/torbenschinke/picd/internal/pic"
	"golang.org/x/sync/semaphore"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type CameraRepo struct {
	mutex  *semaphore.Weighted
	exec   string
	libcam bool
}

func NewCameraRepo() *CameraRepo {
	r := &CameraRepo{mutex: semaphore.NewWeighted(1)}
	r.exec = "raspistill"
	if _, err := os.Stat("/usr/bin/libcamera-still"); err == nil {
		log.Printf("detected libcamera-still instead raspistill\n")
		r.exec = "libcamera-still"
		r.libcam = true
	}
	return r
}

func (r *CameraRepo) FindAllStillCameras() ([]pic.StillCamera, error) {
	panic("implement me")
}

func (r *CameraRepo) CapturePhoto(settings pic.Settings, dst []byte) ([]byte, error) {
	if !r.mutex.TryAcquire(1) {
		return nil, fmt.Errorf("already busy taking a capture")
	}

	defer r.mutex.Release(1)

	var args []string
	if r.libcam {
		// don't have iso yet
		v, _ := strconv.Atoi(string(settings.ISO))
		gainLike := v / 100 // iso-100 => gain=1, iso-200 => gain=2 etc.
		if gainLike > 0 {
			args = append(args, "--gain", strconv.Itoa(gainLike))
		}

	} else {
		if settings.ISO != "" {
			args = append(args, "--ISO", string(settings.ISO))
		}
	}

	if settings.Quality != 0 {
		args = append(args, "--quality", strconv.Itoa(settings.Quality))
	}

	if settings.Resolution.X != 0 {
		args = append(args, "--width", strconv.Itoa(settings.Resolution.X))
	}

	if settings.Resolution.Y != 0 {
		args = append(args, "--height", strconv.Itoa(settings.Resolution.Y))
	}

	if settings.Shutter != 0 {
		us := settings.Shutter.Microseconds()
		args = append(args, "--shutter", strconv.FormatInt(us, 10))
	}

	if settings.AWB != "" {
		args = append(args, "--awb", settings.AWB)
	}

	if settings.Saturation != "" {
		args = append(args, "--saturation", settings.Saturation)
	}

	if settings.Denoise != "" {
		args = append(args, "--denoise", settings.Denoise)
	}

	if settings.EV != "" {
		args = append(args, "--ev", settings.EV)
	}

	if r.libcam {
		var exifRot = "0"
		switch settings.Rotation {
		case 90:
			exifRot = "6"
		case 180:
			exifRot = "3"
		case 270:
			exifRot = "8"
		}
		args = append(args, "--exif", "IFD0.Orientation="+exifRot)
	} else {
		if settings.Rotation != 0 {
			args = append(args, "--rotation", strconv.Itoa(settings.Rotation))
		}
	}

	if !r.libcam {
		// just ignore that, just sport or normal, unusable for us
		if settings.Exposure != "" {
			args = append(args, "--exposure", settings.Exposure)
		}
	}

	if settings.Mode != "" {
		args = append(args, "--mode", settings.Mode)
	}

	if r.libcam {
		args = append(args, "--nopreview")
	}

	args = append(args, "-o", "-") // stream into stdout

	cmd := exec.Command(r.exec, args...)
	cmd.Env = os.Environ()

	fmt.Println(cmd.String())

	errBuf := &bytes.Buffer{}
	imgBuf := bytes.NewBuffer(dst[:0])

	cmd.Stderr = errBuf
	cmd.Stdout = imgBuf

	if err := cmd.Run(); err != nil {
		fmt.Println(errBuf.String())
		return nil, err
	}

	return imgBuf.Bytes(), nil
}
