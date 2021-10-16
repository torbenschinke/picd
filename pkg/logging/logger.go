package logging

import (
	"context"
	"github.com/golangee/log"
	"github.com/golangee/log/field"
)

type Logger = log.Logger

func FromContext(ctx context.Context) Logger {
	return log.FromContext(ctx)
}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return log.WithLogger(ctx, l)
}

func NewLogger(fields ...interface{}) Logger {
	return log.NewLogger()
}

func V(k string, v interface{}) field.DefaultField {
	return log.V(k, v)
}
