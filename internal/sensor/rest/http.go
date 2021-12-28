package rest

import (
	"encoding/json"
	"github.com/torbenschinke/picd/internal/sensor"
	"github.com/torbenschinke/picd/pkg/logging"
	"github.com/torbenschinke/picd/pkg/server"
	"net/http"
)

type Controller struct {
	service *sensor.SenseService
}

func NewController(service *sensor.SenseService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) sensors(w http.ResponseWriter, r *http.Request) {
	type sensors struct {
		Temperature float64
		Humidity    float64
	}

	if err := json.NewEncoder(w).Encode(&sensors{
		Temperature: c.service.T().Celsius(),
		Humidity:    c.service.RH().Humidity(),
	}); err != nil {
		logger := logging.FromContext(r.Context())
		logger.Println(err)
	}
}

func (c *Controller) Configure(r server.Router) {
	r.HandlerFunc(http.MethodGet, "/api/v1/sensors", c.sensors)
}
