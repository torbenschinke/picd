package main

import (
	"context"
	"github.com/torbenschinke/picd/internal/vid"
	"github.com/torbenschinke/picd/pkg/logging"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := logging.NewLogger(
		logging.V("build_id", "tbd"),
		logging.V("build_tag", "tbd"),
	)

	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Println("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Println(err)
	}

	logger.Println("successful shutdown")
}

// TODO http://cat-cam-01:8080/api/v1/camera/1/capture?rotation=180&quality=80&saturation=0&iso=6400&denoise=cdn_hq&ev=4

func realMain(ctx context.Context) error {
	a := vid.NewApplication(ctx)

	return a.Run(ctx)
}
