package midlewares

import (
	"context"
)

type Verifier interface {
	Verify(ctx context.Context, token string) (bool, error)
}

type Midlewares struct {
	ctx context.Context
	v   Verifier
}

func NewMidelware(ctx context.Context, v Verifier) *Midlewares {
	return &Midlewares{
		ctx: ctx,
		v:   v,
	}
}
