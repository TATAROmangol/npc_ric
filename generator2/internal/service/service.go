package service

import (
	"context"
	"fmt"
	"generator/internal/entity"
	"generator/pkg/logger"
	"mime/multipart"
	"os"
)

type Repo interface{
	DeleteTemplate(ctx context.Context, id int) error
	UploadTemplate(ctx context.Context, id int, data []byte) error
	GenerateTemplate(ctx context.Context, id int) ([]byte ,error)
}

type Generatorer interface{
	Generate(ctx context.Context, data []byte, table entity.Table) ([]byte, error)
}

type Tabler interface{
	GetTable(ctx context.Context, institutionId int) (entity.Table, error)
}

type Service struct{
	repo Repo
	generator Generatorer
	tabler Tabler
}

func (s *Service) DeleteTemplate(ctx context.Context, id int) error{
	return s.repo.DeleteTemplate(ctx, id)
}

func (s *Service) UploadTemplate(ctx context.Context, id int, file multipart.File) error{
	var data []byte 
	_, err := file.Read(data)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed file -> []bytes", err)
		return err
	}

	return s.repo.UploadTemplate(ctx, id, data)
}

func (s *Service) GenerateTemplate(ctx context.Context, id int) (*os.File, func(), error){
	data, err := s.repo.GenerateTemplate(ctx, id)
	if err != nil{
		return nil, func() {}, err
	}

	table, err := s.tabler.GetTable(ctx, id)
	if err != nil{
		return nil, func() {}, err
	}

	rdata, err := s.generator.Generate(ctx, data, table)
	if err != nil{
		return nil, func() {}, err
	}

	file, err := os.CreateTemp("", fmt.Sprintf("%v.docx", id))
	if err != nil {
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to create temp file", err)
		return nil, func() {}, err
	}
	file.Write(rdata)

	remove := func() {
		file.Close()
		err := os.Remove(file.Name())
		if err != nil{
			logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to remove file", err)
		}
	}

	return file, remove, nil
}