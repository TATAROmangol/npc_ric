package logger

import (
	"context"
	"log/slog"
)

type LogSender interface{
	SendLog(ctx context.Context, message []byte) error
}

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
    if attrs, ok := ctx.Value(LogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return h.Handler.Handle(ctx, r)
}