package service

import (
	"auth/pkg/logger"
	"context"
)

//go:generate mockgen -source=service.go -destination=./mocks/service_mock.go -package=mocks

type Validator interface {
	IsValid(login, password string) bool
}

type JWT interface {
	GenerateToken() (string, error)
	IsAdmin(tokenString string) (bool, error)
}

type Service struct{
	Validator Validator
	Jwt JWT
}

func New(validator Validator, jwt JWT) *Service {
	return &Service{
		Validator: validator,
		Jwt: jwt,
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	ok := s.Validator.IsValid(login, password)
	if !ok {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "invalid login or password", nil)
		return "", nil
	}

	token, err := s.Jwt.GenerateToken()
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to generate token", err)
		return "", err
	}

	return token, nil
}

func (s *Service) IsAdmin(ctx context.Context,tokenString string) (bool, error){
	ok, err := s.Jwt.IsAdmin(tokenString)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to parse token", err)
		return false, err
	}

	return ok, nil
}

