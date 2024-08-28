package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/cockroachdb/errors"
)

const (
	BatchName = "batch-name"
	RequestId = "request-id"
)

var (
	ErrContextKeyNotFound = errors.New("context key not found")
)

func newLoggerWithContext(ctx context.Context, key string) (*slog.Logger, error) {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return nil, errors.Wrapf(ErrContextKeyNotFound, "key: %s", key)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, NewOption())).With(slog.String(key, value))
	return logger, nil
}

func NewCLILogger(ctx context.Context) (*slog.Logger, error) {
	return newLoggerWithContext(ctx, BatchName)
}

func NewServerLogger(ctx context.Context) (*slog.Logger, error) {
	return newLoggerWithContext(ctx, RequestId)
}

func NewOption() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: true,
		Level:     nil,
	}
}
