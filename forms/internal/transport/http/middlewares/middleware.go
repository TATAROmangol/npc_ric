package middlewares

import (
	"context"
	"forms/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

//go:generate mockgen -source=middleware.go -destination=mocks/middlewares.go -package=mocks

type Verifier interface {
	Verify(ctx context.Context, token string) (bool, error)
}

type Middlewares struct {
	Verifier Verifier
}

func NewMiddlewares(verifier Verifier) *Middlewares {
	return &Middlewares{
		Verifier: verifier,
	}
}

func (m *Middlewares) InitLoggerContextMiddleware(ctx context.Context) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx = logger.SwapContext(ctx, r.Context())
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func (m *Middlewares) InitJsonContentTypeMiddleware() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	}
}

func (m *Middlewares) OperationMiddleware() func(h http.Handler) http.Handler {
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

func (m *Middlewares) AuthMiddleware(cookieName string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("admin_token")
			if err != nil || cookie == nil {	
				logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "cookie not found", err)
				http.Error(w, "cookie not found", http.StatusUnauthorized)
				return
			}

			ok, err := m.Verifier.Verify(r.Context(), cookie.Value)
			if err != nil || !ok {
				logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid token", err)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}