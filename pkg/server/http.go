package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/torbenschinke/picd/pkg/logging"
	"net"
	"net/http"
	"time"
)

type Params = httprouter.Params

type Router interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

func ParamsFromContext(ctx context.Context) Params {
	return httprouter.ParamsFromContext(ctx)
}

type RouterConfigurator interface {
	Configure(Router)
}

func Serve(ctx context.Context, configurators ...RouterConfigurator) error {
	router := httprouter.New()
	for _, configurator := range configurators {
		configurator.Configure(router)
	}

	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      router,
	}

	logger := logging.FromContext(ctx)
	logger.Println("starting server...")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("cannot listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

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
