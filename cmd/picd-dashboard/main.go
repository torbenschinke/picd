package main

import (
	"context"
	"github.com/torbenschinke/picd/internal/dashboard/app"
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

func realMain(ctx context.Context) error {
	a := app.NewApplication()

	return a.Run(ctx)
}
