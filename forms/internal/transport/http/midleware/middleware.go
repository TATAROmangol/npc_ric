package middleware

import (
	"context"
	"forms/pkg/logger"
	"net/http"
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

func InitJsonContentTypeMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	}
}