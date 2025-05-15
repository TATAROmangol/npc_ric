package put

import (
	"context"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -source=put.go -destination=./mocks/put_mock.go -package=mocks

type Putter interface {
	PutInstitutionInfo(ctx context.Context, id int, name string, inn int) error
	PutInstitutionColumns(ctx context.Context, id int, columns []string) error
	PutMentor(ctx context.Context, id int, name string) error
}

type PutHandler struct {
	srv Putter
}

func NewHandler(srv Putter) *PutHandler {
	return &PutHandler{
		srv: srv,
	}
}

func (h *PutHandler) PutInstitutionInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PutInstitutionInfoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.srv.PutInstitutionInfo(r.Context(), req.Id, req.Name, req.INN)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (h *PutHandler) PutInstitutionColumns() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PutInstitutionColumnsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.srv.PutInstitutionColumns(r.Context(), req.Id, req.Columns)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (h *PutHandler) PutMentor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PutMentorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.srv.PutMentor(r.Context(), req.Id, req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
