package delete

import (
	"context"
	"encoding/json"
	"errors"
	"forms/pkg/logger"
	"net/http"
)

//go:generate mockgen -source=delete.go -destination=./mocks/mock_delete.go -package=delete

type Deleter interface {
	DeleteInstitution(ctx context.Context, institutionId int) error
	DeleteMentor(ctx context.Context, mentorId int) error
}

var (
	ErrInvalidInstitutionId = errors.New("invalid institution id")
	ErrInvalidMentorId      = errors.New("invalid mentor id")
)

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
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Id <= 0 {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid institution id", ErrInvalidInstitutionId)
			http.Error(w, ErrInvalidInstitutionId.Error(), http.StatusBadRequest)
			return
		}

		if err := h.srv.DeleteInstitution(r.Context(), req.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (h *DeleteHandler) DeleteMentor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req DeleteMentorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Id <= 0 {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid mentor id", ErrInvalidMentorId)
			http.Error(w, ErrInvalidMentorId.Error(), http.StatusBadRequest)
			return
		}

		if err := h.srv.DeleteMentor(r.Context(), req.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
