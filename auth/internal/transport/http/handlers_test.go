package httpserver

import (
	"auth/internal/transport/http/mocks"
	"auth/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLoginHandler(t *testing.T) {
	loginer := mocks.NewMockLoginer(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(login, password string)

	tests := []struct {
		name string
		mockBehavior MockBehavior
		req LoginRequest
		want string
		wantStatus int
	}{
		{
			name: "valid credentials",
			want: "test",
			wantStatus: http.StatusOK,
			req: LoginRequest{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(login, password string) {
				loginer.EXPECT().Login(ctx, login, password).Return("test", nil)
			},
		},
		{
			name: "wrong password",
			want: "",
			wantStatus: http.StatusBadRequest,
			req: LoginRequest{
				Login:    "login",
				Password: "",
			},
			mockBehavior: func(login, password string) {},
		},
		{
			name: "wrong password",
			want: "",
			wantStatus: http.StatusBadRequest,
			req: LoginRequest{
				Login:    "",
				Password: "password",
			},
			mockBehavior: func(login, password string) {},
		},
		{
			name: "err token",
			want: "",
			wantStatus: http.StatusInternalServerError,
			req: LoginRequest{
				Login:    "login",
				Password: "password",
			},
			mockBehavior: func(login, password string) {
				loginer.EXPECT().Login(ctx, login, password).Return("", errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.req.Login, tt.req.Password)

			rr := httptest.NewRecorder()

			json, _ := json.Marshal(tt.req)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(json))
			req.Header.Set("Content-Type", "application/json")
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			LoginHandler(loginer).ServeHTTP(rr, req)

			if tt.wantStatus != rr.Code {
				t.Errorf("LoginHandler status got %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			cookie := rr.Result().Cookies()
			if len(cookie) == 0 {
				t.Errorf("LoginHandler cookie got %v, want %v", cookie, 1)
			}
			if cookie[0].Name != "admin_token" {
				t.Errorf("LoginHandler cookie name got %v, want %v", cookie[0].Name, "admin_token")
			}
			if cookie[0].Value != tt.want {
				t.Errorf("LoginHandler cookie value got %v, want %v", cookie[0].Value, tt.want)
			}
		})
	}
}
