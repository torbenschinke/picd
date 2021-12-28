package app

import (
	"context"
	"github.com/torbenschinke/picd/internal/iam/none"
	"github.com/torbenschinke/picd/internal/pic"
	"github.com/torbenschinke/picd/internal/pic/raspistill"
	"github.com/torbenschinke/picd/internal/pic/rest"
	"github.com/torbenschinke/picd/internal/sensor"
	rest2 "github.com/torbenschinke/picd/internal/sensor/rest"
	"github.com/torbenschinke/picd/pkg/server"
)

type Application struct {
}

func NewApplication(ctx context.Context) *Application {
	return &Application{}
}

func (a *Application) Run(ctx context.Context) error {
	camRepo := raspistill.NewCameraRepo()
	picService := pic.NewCameraService(camRepo)
	httpController := rest.NewController(picService, none.Authenticator{})
	senHttpCtr := rest2.NewController(sensor.NewSenseService())

	return server.Serve(ctx, httpController, senHttpCtr)
}
