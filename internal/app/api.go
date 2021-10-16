package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/torbenschinke/picd/internal/iam/none"
	"github.com/torbenschinke/picd/internal/pic"
	"github.com/torbenschinke/picd/internal/pic/raspistill"
	"github.com/torbenschinke/picd/internal/pic/rest"
	"github.com/torbenschinke/picd/pkg/logging"
	"net"
	"net/http"
	"time"
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

	router := httprouter.New()
	rest.Configure(router, httpController)

	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      router,
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("cannot listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		logger := logging.FromContext(ctx)
		logger.Println("Serve: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Println("Serve: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("failed to shutdown: %w", err)
	}

	return nil
}
