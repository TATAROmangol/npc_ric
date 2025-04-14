package handlers

import (
	"net/http"
)

type Creater interface{
}

func (h *Handler) CreateInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (h *Handler) DeleteInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (h *Handler) CreateMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (h *Handler) DeleteMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}