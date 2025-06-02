package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	log *slog.Logger
}

func New() *Logger {
	log := slog.New(
		&ContextHandler{slog.NewJSONHandler(os.Stdout, nil)},
	)
	return &Logger{log}
}

func AppendCtx(ctx context.Context, key string, val any) context.Context {
	attr := slog.Any(key, val)

	if ctx == nil {
		ctx = context.Background()
	}

	if v, ok := ctx.Value(LogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(ctx, LogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(ctx, LogFields, v)
}

func SwapContext(ctxWithLogger, other context.Context) context.Context {
	l := GetFromCtx(ctxWithLogger)

	if v, ok := ctxWithLogger.Value(LogFields).([]slog.Attr); ok {
		other = context.WithValue(other, LogFields, v)
	}

	return InitFromCtx(other, l)
}

func InitFromCtx(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, Key, log)
}

func GetFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log.InfoContext(ctx, msg, args...)
}

func (l Logger) ErrorContext(ctx context.Context, msg string, err error) {
	l.log.ErrorContext(ctx, msg, "error", err)
}
