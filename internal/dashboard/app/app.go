package app

import (
	"context"
	"github.com/torbenschinke/picd/internal/dashboard"
	"github.com/torbenschinke/picd/internal/dashboard/html"
	"github.com/torbenschinke/picd/pkg/server"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Run(ctx context.Context) error {
	dashboardService := dashboard.NewService()
	handler := html.NewHandler(dashboardService)

	return server.Serve(ctx, handler)
}
