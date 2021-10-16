package rest

import (
	"fmt"
	"github.com/torbenschinke/picd/internal/pic"
	"github.com/torbenschinke/picd/pkg/logging"
	"github.com/torbenschinke/picd/pkg/server"
	"net/http"
	"strconv"
)

type Authenticator interface {
	// Authenticate returns true if everything is fine, otherwise returns false and writes
	// the according status code back.
	Authenticate(w http.ResponseWriter, r *http.Request) bool
}

type Controller struct {
	service       *pic.CameraService
	authenticator Authenticator
}

func NewController(service *pic.CameraService, authenticator Authenticator) *Controller {
	return &Controller{service: service, authenticator: authenticator}
}

// captureImage is e.g. GET /api/v1/camera/:id/capture?iso=100&awb=auto&x=1920&y=1080.
func (c *Controller) captureImage(w http.ResponseWriter, r *http.Request) {
	logger := logging.FromContext(r.Context())
	if !c.authenticator.Authenticate(w, r) {
		return
	}

	params := server.ParamsFromContext(r.Context())
	cameraId := pic.CameraID(params.ByName("id"))
	settings := pic.Settings{
		Camera: cameraId,
	}

	query := r.URL.Query()
	iso := query.Get("iso")
	if iso != "" {
		settings.ISO = pic.ISO(iso)
	}

	x, err := strconv.Atoi(query.Get("x"))
	if err != nil {
		logger.Println(fmt.Errorf("invalid x value: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	y, err := strconv.Atoi(query.Get("y"))
	if err != nil {
		logger.Println(fmt.Errorf("invalid y value: %w", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings.Resolution.X = x
	settings.Resolution.Y = y

	img, err := c.service.CapturePhoto(settings)
	if err != nil {
		logger.Println(fmt.Errorf("cannot capture photo: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer img.Recycle()

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = img.WriteTo(w)
	if err != nil {
		logger.Println(fmt.Errorf("cannot write image response: %w", err))
	}
}

func Configure(r server.Router, c *Controller) {
	r.HandlerFunc(http.MethodGet, "/api/v1/camera/:id/capture", c.captureImage)
}
