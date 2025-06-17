package table

import (
	"context"
	"generator/internal/entity"
	tablepb "generator/pkg/grpc/table"
	"generator/pkg/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct{
	cfg Config
	conn *grpc.ClientConn
	client tablepb.TableServiceClient
}

func New() *Client {
	return &Client{}
}

func (c *Client) connect(ctx context.Context) error{
	var err error
	for i := 0; i < 5; i++{
		if c.client != nil{
			break
		}

		con, err := grpc.NewClient(
			c.cfg.Addr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to create table grpc client", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		c.conn = con
		c.client = tablepb.NewTableServiceClient(con)
	}

	if c.client == nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to create table grpc client all attempts", err)
		return err
	}

	logger.GetFromCtx(ctx).InfoContext(ctx, "Listen table", "path", c.cfg.Addr())
	return nil
}

func (c *Client) GetTable(ctx context.Context, institutionId int) (entity.Table, error) {
	if err := c.connect(ctx); err != nil {
		return entity.Table{}, err
	}

	resp, err := c.client.GetTable(ctx, &tablepb.GetTableRequest{
		InstitutionId: int32(institutionId),
	})

	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to get table", err)
		return entity.Table{}, err
	}

	var rows [][]string 
	for _, row := range resp.Rows {
		rows = append(rows, row.Values)
	}

	return entity.Table{Columns: resp.Columns, Rows: rows}, nil
}

func (c *Client) Close() error{
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.GetFromCtx(context.Background()).ErrorContext(context.Background(), "failed to close table grpc connection", err)
			return err
		}
	}

	return nil
}

