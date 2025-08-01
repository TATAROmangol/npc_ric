package tablegrpc

import (
	"context"
	"forms/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func InitLogger(pCtx context.Context) func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (any, error) {

		ctx = logger.InitFromCtx(ctx, logger.GetFromCtx(pCtx))

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func Operation() func(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (any, error) {

		operatiodID := uuid.New()
		ctx = logger.AppendCtx(ctx, "operation_id", operatiodID.String())
		ctx = logger.AppendCtx(ctx, "method_name", info.FullMethod)

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}