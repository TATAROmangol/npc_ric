package httpserver

import (
	"auth/pkg/logger"
	"context"
	"net/http"

	"github.com/google/uuid"
)

func InitLoggerContextMiddleware(ctx context.Context) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx = logger.SwapContext(ctx, r.Context())
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func Operation() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			operationId := uuid.NewString()

			ctx := r.Context()
			ctx = logger.AppendCtx(ctx, "operation_id", operationId)
			ctx = logger.AppendCtx(ctx, "method path", r.URL.Path)
			logger.GetFromCtx(ctx).InfoContext(ctx, "called method")

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}