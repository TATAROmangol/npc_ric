package service

import (
	"context"
	"mime/multipart"
	"os"
)

type Service struct{

}

func (s *Service) DeleteTemplate(ctx context.Context, id int) error{
	return nil
}

func (s *Service) UploadTemplate(id int, file multipart.File) error{
	var data []byte 
	_, err := file.Read(data)
	if err != nil{

	}
	return nil
}

func (s *Service) GenerateTemplate(id int) (*os.File ,error){
	return nil,nil
}