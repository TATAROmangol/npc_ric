package delete

import (
	"context"
	"encoding/json"
	"net/http"
)

type Deleter interface {
	DeleteInstitution(ctx context.Context, institutionId int) error
	DeleteMentor(ctx context.Context, mentorId int) error
}

type DeleteHandler struct {
	srv Deleter
}

func NewHandler(srv Deleter) *DeleteHandler {
	return &DeleteHandler{
		srv: srv,
	}
}

func (h *DeleteHandler) DeleteInstitution() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req DeleteInstitutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		if err := h.srv.DeleteInstitution(r.Context(), req.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *DeleteHandler) DeleteMentor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req DeleteMentorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		if err := h.srv.DeleteMentor(r.Context(), req.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
