package midlewares

import (
	"context"
)

type Verifier interface{
	Verify(ctx context.Context, token string) (bool, error)
}

type Midleware struct {
	ctx context.Context
	v Verifier
}

func NewMidelware(ctx context.Context, v Verifier) *Midleware{
	return &Midleware{
		ctx:ctx,
		v:v,
	}
}
