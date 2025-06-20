package service

import (
	"context"
	"generator/internal/entity"
)

const(
	Dir = "../temp"
)

type Repo interface{
	DeleteTemplate(ctx context.Context, id int) error
	UploadTemplate(ctx context.Context, id int, data []byte) error
	GetTemplate(ctx context.Context, id int) ([]byte ,error)
}

type Generatorer interface{
	Generate(data []byte, table entity.Table) ([]byte, error)
}

type Tabler interface{
	GetTable(ctx context.Context, institutionId int) (entity.Table, error)
}

type Service struct{
	repo Repo
	generator Generatorer
	tabler Tabler
}

func New(repo Repo, generator Generatorer, tabler Tabler) *Service{
	return &Service{
		repo: repo,
		generator: generator,
		tabler: tabler,
	}
}

func (s *Service) DeleteTemplate(ctx context.Context, id int) error{
	return s.repo.DeleteTemplate(ctx, id)
}

func (s *Service) UploadTemplate(ctx context.Context, id int, file []byte) error{
	return s.repo.UploadTemplate(ctx, id, file)
}

func (s *Service) GenerateTemplate(ctx context.Context, id int) ([]byte, error){
	data, err := s.repo.GetTemplate(ctx, id)
	if err != nil{
		return nil, err
	}

	table, err := s.tabler.GetTable(ctx, id)
	if err != nil{
		return nil, err
	}

	file, err := s.generator.Generate(data, table)
	if err != nil{
		return nil, err
	}

	return file, nil
}