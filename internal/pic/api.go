package pic

import (
	v1 "github.com/torbenschinke/picd/pkg/api/v1"
	"io"
	"sync"
)

type ISO = v1.ISO
type CameraID = v1.CameraID
type StillCamera = v1.StillCamera
type Settings = v1.Settings

type Image interface {
	Recycle()
	io.WriterTo
}

// CameraService provides the domains' use cases.
type CameraService struct {
	repo StillCameraRepository
	pool *sync.Pool
}

func NewCameraService(repo StillCameraRepository) *CameraService {
	s := &CameraService{repo: repo}
	s.pool = &sync.Pool{New: func() interface{} {
		return &img{
			pool: s.pool,
		}
	}}

	return s
}

func (s *CameraService) CapturePhoto(settings Settings) (Image, error) {
	image := s.pool.Get().(*img)
	buf := image.buf
	buf, err := s.repo.CapturePhoto(settings, buf)
	if err != nil {
		image.Recycle()
		return nil, err
	}

	image.buf = buf

	return image, nil
}

func (s *CameraService) ListCameras() ([]StillCamera, error) {
	return s.repo.FindAllStillCameras()
}
