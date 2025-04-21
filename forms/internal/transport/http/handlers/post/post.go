package post

import (
	"context"
	"encoding/json"
	"forms/pkg/logger"
	"net/http"
)

type Poster interface {
	PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error)
	PostForm(ctx context.Context, institutionId int, info []string) (int, error)
}

type PostHandler struct {
	srv Poster
}

func NewHandler(srv Poster) *PostHandler {
	return &PostHandler{
		srv: srv,
	}
}

func (h *PostHandler) PostInstitution() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PostInstitutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := h.srv.PostInstitution(r.Context(), req.Name, req.INN, req.Columns)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(PostInstitutionResponse{Id: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *PostHandler) PostForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PostFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		id, err := h.srv.PostForm(r.Context(), req.InstitutionId, req.Info)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(PostFormResponse{Id: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
