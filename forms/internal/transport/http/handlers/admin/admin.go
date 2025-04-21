package admin

import (
	"net/http"
)

type Creater interface{
}

type AdminHandler struct{
	srv Creater
}

func NewHandler(srv Creater) *AdminHandler {
	return &AdminHandler{
		srv: srv,
	}
}

func (h *AdminHandler) CreateInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *AdminHandler) TransformInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *AdminHandler) DeleteInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *AdminHandler) CreateMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *AdminHandler) TransformMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (h *AdminHandler) DeleteMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}