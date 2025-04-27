package repository

import (
	"database/sql"
	"forms/internal/storage/repository/delete"
	"forms/internal/storage/repository/get"
	"forms/internal/storage/repository/post"
	"forms/internal/storage/repository/put"
)


type Storage struct {
	*get.Get
	*put.Put
	*delete.Delete
	*post.Post
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Get: get.NewGet(db),
		Put: put.NewPut(db),
		Delete: delete.NewDelete(db),
		Post: post.NewPost(db),
	}
}