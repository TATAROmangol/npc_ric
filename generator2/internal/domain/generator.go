package domain

import (
	"context"
	"fmt"
	"generator/internal/entity"
)

type Generator struct{}

func NewGenerator() *Generator {
	 return &Generator{}
}

func (g *Generator) Generate(ctx context.Context, data []byte, table entity.Table) ([]byte, error){
	return nil, fmt.Errorf("not implemented")
}