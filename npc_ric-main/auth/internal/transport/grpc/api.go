package grpcserver

import (
	ssov1 "auth/pkg/grpc/auth"
	"context"

	"google.golang.org/grpc"
)

type Verificate interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

type Api struct{
	ssov1.UnimplementedVerifyServer
	Verificate
}

func Register(gRPCServer *grpc.Server, ver Verificate) {
	ssov1.RegisterVerifyServer(gRPCServer, &Api{Verificate: ver})
}

func (a *Api) Verify(ctx context.Context, req *ssov1.VerifyRequest) (*ssov1.VerifyResponse, error){
	ok, err := a.Verificate.IsAdmin(ctx, req.Token)
	if err != nil{
		return nil, err
	}

	return &ssov1.VerifyResponse{
		IsAdmin: ok,
	}, nil
}