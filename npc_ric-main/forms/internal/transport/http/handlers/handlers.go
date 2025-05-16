package handlers

import (
	"forms/internal/transport/http/handlers/delete"
	"forms/internal/transport/http/handlers/get"
	"forms/internal/transport/http/handlers/post"
	"forms/internal/transport/http/handlers/put"
	"net/http"
)


type Service interface {
	delete.Deleter
	get.Getter
	put.Putter
	post.Poster
}

type Handlers struct{
	Deleter *delete.DeleteHandler
	Getter *get.GetHandler
	Putter *put.PutHandler
	Poster *post.PostHandler
}

func NewHandlers(srv Service) *Handlers {
	return &Handlers{
		Deleter: delete.NewHandler(srv),
		Getter: get.NewHandler(srv),
		Putter: put.NewHandler(srv),
		Poster: post.NewHandler(srv),
	}
}

func (h *Handlers) DeleteInstitution() http.Handler {
	return h.Deleter.DeleteInstitution()
}

func (h *Handlers) DeleteMentor() http.Handler {
	return h.Deleter.DeleteMentor()
}

func (h *Handlers) GetInstitutions() http.Handler {
	return h.Getter.GetInstitutions()
}
func (h *Handlers) GetMentors() http.Handler {
	return h.Getter.GetMentors()
}
func (h *Handlers) GetInstitutionFromINN() http.Handler {
	return h.Getter.GetInstitutionFromINN()
}
func (h *Handlers) GetFormColumns() http.Handler {
	return h.Getter.GetFormColumns()
}
func (h *Handlers) PutInstitutionInfo() http.Handler {
	return h.Putter.PutInstitutionInfo()
}
func (h *Handlers) PutInstitutionColumns() http.Handler {
	return h.Putter.PutInstitutionColumns()
}
func (h *Handlers) PutMentor() http.Handler {
	return h.Putter.PutMentor()
}
func (h *Handlers) PostInstitution() http.Handler {
	return h.Poster.PostInstitution()
}
func (h *Handlers) PostMentor() http.Handler {
	return h.Poster.PostMentor()
}
func (h *Handlers) PostForm() http.Handler {
	return h.Poster.PostForm()
}