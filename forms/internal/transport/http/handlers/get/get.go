package get

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
	GetInstitutionFromINN(ctx context.Context, inn int) (entities.Institution, error)
	GetFormColumns(ctx context.Context, id int) ([]string, error)
}

type GetHandler struct {
	srv Getter
}

func NewHandler(srv Getter) *GetHandler {
	return &GetHandler{
		srv: srv,
	}
}

func (h *GetHandler) GetInstitutions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		institutions, err := h.srv.GetInstitutions(r.Context())
		if err != nil {
			http.Error(w, "failed to get institutions id", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(institutions); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *GetHandler) GetMentors() http.Handler {
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

func (h *GetHandler) GetInstitutionFromINN() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req GetInstitutionFromINNRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		institution, err := h.srv.GetInstitutionFromINN(r.Context(), req.Inn)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(institution); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *GetHandler) GetFormColumns() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req GetFormColumnsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		columns, err := h.srv.GetFormColumns(r.Context(), req.InstitutionId)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if len(columns) == -1 {
			http.Error(w, "invalid request", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(columns); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
