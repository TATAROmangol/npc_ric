package service

import (
	"context"
	"fmt"
	"generator/internal/entity"
	"generator/pkg/logger"
	"mime/multipart"
	"os"
	"time"
)

const(
	Dir = "./temp"
)

type Repo interface{
	DeleteTemplate(ctx context.Context, id int) error
	UploadTemplate(ctx context.Context, id int, data []byte) error
	GetTemplate(ctx context.Context, id int) ([]byte ,error)
}

type Generatorer interface{
	Generate(data []byte, table entity.Table, path string) error 
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

func (s *Service) UploadTemplate(ctx context.Context, id int, file multipart.File) error{
	var data []byte 
	_, err := file.Read(data)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed file -> []bytes", err)
		return err
	}

	return s.repo.UploadTemplate(ctx, id, data)
}

func (s *Service) GenerateTemplate(ctx context.Context, id int) ([]byte, func(), error){
	data, err := s.repo.GetTemplate(ctx, id)
	if err != nil{
		return nil, func() {}, err
	}

	table, err := s.tabler.GetTable(ctx, id)
	if err != nil{
		return nil, func() {}, err
	}

	path := fmt.Sprintf("%s/%v.docx", Dir, id)
	err = s.generator.Generate(data, table, path)
	if err != nil{
		return nil, func() {}, err
	}

	file, err := os.ReadFile(path)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to read file", err)
		return nil, func() {}, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	go func() {
		<-ctx.Done()
		
		if err := os.Remove(path); err != nil {
			logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to remove file", err)
		}
	}()

	remove := func() {
		cancel() 
	}

	return file, remove, nil
}