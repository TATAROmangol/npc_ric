package service

import (
	"forms/internal/service/get"
	"forms/internal/service/post"
	"forms/internal/service/put"
	"forms/internal/service/delete"
)

type Repo interface {
	get.GetRepo
	put.PutRepo
	post.PostRepo
	delete.DeleteRepo
}

type Services struct {
	*get.GetService
	*put.PutService
	*post.PostService
	*delete.DeleteService
}

func NewServices(repo Repo) *Services {
	return &Services{
		GetService:  get.NewGetService(repo),
		PutService:  put.NewPutService(repo),
		PostService: post.NewPostService(repo),
		DeleteService: delete.NewDeleteService(repo),
	}
}