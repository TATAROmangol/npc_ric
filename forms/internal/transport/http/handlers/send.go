package handlers

import (
	"context"
	"encoding/json"
	"forms/pkg/logger"
	"net/http"
)

type Sender interface{
	SendForm(ctx context.Context, institution string, info []string) (int, error)
}

func (h *Handler) SendForm() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SendFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to decode request", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		id, err := h.srv.SendForm(r.Context(), req.Institution, req.Info)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(SendFormResponse{ID: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to encode response", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
