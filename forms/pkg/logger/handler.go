package logger

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
    initFields(ctx, r)

	return h.Handler.Handle(ctx, r)
}

func initFields(ctx context.Context, r slog.Record) {
    if attrs, ok := ctx.Value(LogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
}