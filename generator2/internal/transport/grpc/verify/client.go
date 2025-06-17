package verify

import (
	"context"
	authpb "generator/pkg/grpc/auth"
	"generator/pkg/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg Config
	conn *grpc.ClientConn
	client authpb.VerifyClient
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		conn: nil,
		client: nil,
	}
}

func (c *Client) connect(ctx context.Context) error {
	var err error
	for i := 0; i < 5; i++ {
		if c.client != nil {
			break
		}

		con, err := grpc.NewClient(
			c.cfg.Addr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to create auth grpc client", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		c.conn = con
		c.client = authpb.NewVerifyClient(con)
	}

	if c.conn == nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to createauth grpc client all attempts", err)
		return err
	}

	logger.GetFromCtx(ctx).InfoContext(ctx, "Listen auth", "path", c.cfg.Addr())
	return nil
}

func (c *Client) Verify(ctx context.Context, token string) (bool, error) {
	if c.client == nil {
		if err := c.connect(ctx); err != nil {
			return false, err
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	res, err := c.client.Verify(ctx, &authpb.VerifyRequest{
		Token: token,
	})
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to verify token", err)
		return false, err
	}

	return res.GetIsAdmin(), nil
}

func (c *Client) Close() error{
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.GetFromCtx(context.Background()).ErrorContext(context.Background(), "failed to close auth grpc connection", err)
			return err
		}
	}

	return nil
}