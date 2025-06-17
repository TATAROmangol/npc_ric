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
			cookie, err := r.Cookie("admin_token")
			if err != nil || cookie == nil {	
				logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "cookie not found", err)
				http.Error(w, "cookie not found", http.StatusUnauthorized)
				return
			}

			ok, err := m.v.Verify(r.Context(), cookie.Value)
			if err != nil || !ok {
				logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid token", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}