package httpcreate

import (
	"context"
	"encoding/json"
	"forms/internal/entities"
	"net/http"
)

type Creater interface{
	GetInstitutionsId(ctx context.Context) ([]entities.Institution, error)
}

type FromsHandler struct{
	Creater Creater
}

func (fh *FromsHandler) GetInstitutions() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		institutions, err := fh.Creater.GetInstitutionsId(r.Context())
		if err != nil {
			http.Error(w, "failed to get institutions id", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(GetInstitutionsResponse{Institutions: institutions}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (fh *FromsHandler) CreateInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (fh *FromsHandler) RemoveInstitution() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (fh *FromsHandler) GetMentors() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (fh *FromsHandler) CreateMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}

func (fh *FromsHandler) DeleteMentor() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}