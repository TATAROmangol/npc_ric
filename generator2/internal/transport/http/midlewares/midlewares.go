package midlewares

import (
	"context"
	"generator/pkg/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Verifier interface {
	Verify(ctx context.Context, token string) (bool, error)
}

type Midlewares struct {
	ctx context.Context
	v   Verifier
}

func NewMidelware(ctx context.Context, v Verifier) *Midlewares {
	return &Midlewares{
		ctx: ctx,
		v:   v,
	}
}

func (m *Midlewares) InitLoggerCtx() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.SwapContext(m.ctx, r.Context())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Midlewares) InitJsonContentType() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	}
}

func (m *Midlewares) Operation() func(h http.Handler) http.Handler {
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

func (m *Midlewares) CheckAuth(cookieName string) mux.MiddlewareFunc {
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
