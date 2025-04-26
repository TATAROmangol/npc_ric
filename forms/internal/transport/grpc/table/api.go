package tablegrpc

import (
	"context"
	"fmt"
	tablepb "forms/pkg/grpc/table"
	"forms/pkg/logger"

	"google.golang.org/grpc"
)

type Former interface {
	GetFormRows(ctx context.Context, id int) ([]string,  error)
	GetFormColumns(ctx context.Context, id int) ([]string, error)
}

type Api struct {
	tablepb.UnimplementedTableServiceServer
	Former
}

func Register(gRPCServer *grpc.Server, tb Former) {
	tablepb.RegisterTableServiceServer(gRPCServer, &Api{Former: tb})
}

func (a *Api) GetTable(ctx context.Context, req *tablepb.GetTableRequest) (*tablepb.GetTableResponse, error){
	id := int(req.GetInstitutionId())
	if id <= 0 {
		err := fmt.Errorf("invalid institution id: %d", id)
		logger.GetFromCtx(ctx).ErrorContext(ctx, "invalid institution id", err)
		return nil, err
	}

	rows, err := a.Former.GetFormRows(ctx, id)
	if err != nil {
		return nil, err
	}

	columns, err := a.Former.GetFormColumns(ctx, id)
	if err != nil {
		return nil, err
	}

	return &tablepb.GetTableResponse{
		Rows:    rows,
		Columns: columns,
	}, nil
}






