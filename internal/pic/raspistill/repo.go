package raspistill

import (
	"bytes"
	"fmt"
	"github.com/torbenschinke/picd/internal/pic"
	"golang.org/x/sync/semaphore"
	"os"
	"os/exec"
	"strconv"
)

type CameraRepo struct {
	mutex *semaphore.Weighted
}

func NewCameraRepo() *CameraRepo {
	r := &CameraRepo{mutex: semaphore.NewWeighted(1)}
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
	if settings.ISO != "" {
		args = append(args, "--ISO", string(settings.ISO))
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

	args = append(args, "-o", "-") // stream into stdout

	cmd := exec.Command("raspistill", args...)
	cmd.Env = os.Environ()

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
