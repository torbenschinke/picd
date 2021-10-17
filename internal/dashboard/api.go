package dashboard

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type CameraNode struct {
	ID       string
	Name     string
	Host     string
	Rotation int
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListCameras() ([]CameraNode, error) {
	return []CameraNode{
		{
			ID:       "2",
			Name:     "Stube oben",
			Host:     "http://bahnsen-sat-2:8080",
			Rotation: 180,
		},
		{
			ID:       "1",
			Name:     "Garten",
			Host:     "http://bahnsen-sat-1:8080",
			Rotation: 0,
		},
	}, nil
}

func (s *Service) LoadCameraImage(id string, params url.Values, w io.Writer) error {
	cams, err := s.ListCameras()
	if err != nil {
		return fmt.Errorf("cannot list cameras: %w", err)
	}

	var camNode CameraNode
	for _, cam := range cams {
		if cam.ID == id {
			camNode = cam
			break
		}
	}

	if camNode.ID != id {
		return fmt.Errorf("camera unknown")
	}

	resp, err := http.Get(camNode.Host + "/api/v1/camera/1/capture?" + params.Encode())
	if err != nil {
		return fmt.Errorf("cannot reach camera node: %w", err)
	}

	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("cannot proxy camera image: %w", err)
	}

	return nil
}
