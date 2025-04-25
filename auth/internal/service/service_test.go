package service

import (
	"auth/internal/service/mocks"
	"auth/pkg/logger"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestService_Login(t *testing.T) {
	validator := mocks.NewMockValidator(gomock.NewController(t))
	jwt := mocks.NewMockJWT(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type args struct {
		login    string
		password string
	}

	type MockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         string
		wantErr      bool
	}{
		{
			name: "valid credentials",
			args: args{
				login:    "admin",
				password: "password",
			},
			mockBehavior: func(args args) {
				validator.EXPECT().IsValid(args.login, args.password).Return(true)
				jwt.EXPECT().GenerateToken().Return("token", nil)
			},
			want:    "token",
			wantErr: false,
		},
		{
			name: "wrong password",
			args: args{
				login:    "admin",
				password: "wrongpassword",
			},
			mockBehavior: func(args args) {
				validator.EXPECT().IsValid(args.login, args.password).Return(false)
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "wrong login",
			args: args{
				login:    "wrongadmin",
				password: "password",
			},
			mockBehavior: func(args args) {
				validator.EXPECT().IsValid(args.login, args.password).Return(false)
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "empty valid",
			args: args{
				login:    "",
				password: "",
			},
			mockBehavior: func(args args) {
				validator.EXPECT().IsValid(args.login, args.password).Return(false)
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "token generation error",
			args: args{
				login:    "admin",
				password: "password",
			},
			mockBehavior: func(args args) {
				validator.EXPECT().IsValid(args.login, args.password).Return(true)
				jwt.EXPECT().GenerateToken().Return("", errors.New("test"))
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				validator: validator,
				jwt:       jwt,
			}
			tt.mockBehavior(tt.args)
			got, err := s.Login(ctx, tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_IsAdmin(t *testing.T) {
	validator := mocks.NewMockValidator(gomock.NewController(t))
	jwt := mocks.NewMockJWT(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(token string)

	tests := []struct {
		name    string
		token  string
		mockBehavior MockBehavior
		want    bool
		wantErr bool
	}{
		{
			name: "valid token",
			token: "valid",
			mockBehavior: func(token string) {
				jwt.EXPECT().IsAdmin(token).Return(true, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid token",
			token: "not valid",
			mockBehavior: func(token string) {
				jwt.EXPECT().IsAdmin(token).Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "err token",
			token: "not valid",
			mockBehavior: func(token string) {
				jwt.EXPECT().IsAdmin(token).Return(false, errors.New("test"))
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				validator: validator,
				jwt:       jwt,
			}
			tt.mockBehavior(tt.token)

			got, err := s.IsAdmin(ctx, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.IsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}
