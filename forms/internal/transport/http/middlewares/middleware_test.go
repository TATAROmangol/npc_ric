package middlewares

import (
	"context"
	"forms/internal/transport/http/middlewares/mocks"
	"forms/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMiddlewares_InitLoggerContextMiddleware(t *testing.T) {
	verifier := mocks.NewMockVerifier(gomock.NewController(t))
	m := NewMiddlewares(verifier)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.GetFromCtx(r.Context())
		if logger == nil {
			t.Error("Logger not found in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	l := logger.New()
	ctx := logger.InitFromCtx(context.Background(), l)

	middleware := m.InitLoggerContextMiddleware(ctx)
	middleware(handler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestMiddlewares_InitJsonContentTypeMiddleware(t *testing.T) {
	verifier := mocks.NewMockVerifier(gomock.NewController(t))
	m := NewMiddlewares(verifier)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	middleware := m.InitJsonContentTypeMiddleware()
	middleware(handler).ServeHTTP(rr, req)

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}
}

func TestMiddlewares_AuthMiddleware(t *testing.T) {
	verifier := mocks.NewMockVerifier(gomock.NewController(t))
	m := NewMiddlewares(verifier)

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(val string)

	tests := []struct {
		name   string
		MockBehavior
		cookie http.Cookie
		want   int
	}{
		{
			name: "valid cookie",
			MockBehavior: func(val string) {
				verifier.EXPECT().Verify(gomock.Any(), val).Return(true, nil)
			},
			cookie: http.Cookie{
				Name:  "admin_token",
				Value: "valid_token",
			},
			want:   http.StatusOK,
		},
		{
			name: "invalid cookie",
			MockBehavior: func(val string) {
				verifier.EXPECT().Verify(gomock.Any(), val).Return(false, nil)
			},
			cookie: http.Cookie{
				Name:  "admin_token",
				Value: "invalid_token",
			},
			want:   http.StatusUnauthorized,
		},
		{
			name: "missing cookie",
			MockBehavior: func(val string) {},
			cookie: http.Cookie{},
			want:   http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.cookie.Value)

			req := httptest.NewRequest("GET", "/test", nil)

    		rr := httptest.NewRecorder()

			testHandler := func() func(h http.Handler) http.Handler {
				return func(h http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						r.AddCookie(&tt.cookie)
						h.ServeHTTP(w, r)
					})
				}
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := m.InitLoggerContextMiddleware(ctx)((testHandler()(m.AuthMiddleware()(h))))
			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("Expected status %d, got %d", tt.want, rr.Code)
			}
		})
	}
}
