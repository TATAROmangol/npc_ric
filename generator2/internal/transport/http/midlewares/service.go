package midlewares

import (
	"generator/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *Midleware) InitLoggerCtx() mux.MiddlewareFunc{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.SwapContext(m.ctx, r.Context())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Midleware) CheckAuth(cookieName string) mux.MiddlewareFunc{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}