package handlers

import (
	"context"
	"encoding/json"
	"forms/internal/entities"
	"forms/pkg/logger"
	"net/http"
)

type Getter interface {
	GetInstitutions(ctx context.Context) ([]entities.Institution, error)
	GetMentors(ctx context.Context) ([]entities.Mentor, error)
}

func (h *Handler) GetInstitutions() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		institutions, err := h.srv.GetInstitutions(r.Context())
		if err != nil {
			http.Error(w, "failed to get institutions id", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(institutions); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (h *Handler) GetMentors() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mentors, err := h.srv.GetMentors(r.Context())
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mentors); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *Handler) GetInstitutionFromINN() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


	})
}