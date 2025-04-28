package tablegrpc

import (
	"context"
	"fmt"
	tablepb "forms/pkg/grpc/table"
	"forms/pkg/logger"

	"google.golang.org/grpc"
)

type Former interface {
	GetFormRows(ctx context.Context, id int) ([][]string,  error)
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

	grpcRows := make([]*tablepb.Row, 0, len(rows))
	for _, row := range rows {
		grpcRow := &tablepb.Row{
			Values: row,
		}
		grpcRows = append(grpcRows, grpcRow)
	}

	columns, err := a.Former.GetFormColumns(ctx, id)
	if err != nil {
		return nil, err
	}

	return &tablepb.GetTableResponse{
		Rows:    grpcRows,
		Columns: columns,
	}, nil
}






