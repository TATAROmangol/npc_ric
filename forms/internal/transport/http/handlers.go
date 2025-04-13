package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
)

type Former interface{
	SendForm(ctx context.Context, institution string, info []string) (int, error)
	GetInstitutionsId(ctx context.Context) ([]int, error)
}

type FromsHandler struct{
	Former Former
}

func (fh *FromsHandler) SendForm() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct{
			Institution string   `json:"institution"`
			Info        []string `json:"info"`
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		id, err := fh.Former.SendForm(r.Context(), req.Institution, req.Info)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		type Response struct{
			ID int `json:"id"`
		}

		if err := json.NewEncoder(w).Encode(Response{ID: id}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (fh *FromsHandler) GetInstitutionsId() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := fh.Former.GetInstitutionsId(r.Context())
		if err != nil {
			http.Error(w, "failed to get institutions id", http.StatusInternalServerError)
			return
		}

		type Response struct{
			ID []int `json:"id"`
		}

		if err := json.NewEncoder(w).Encode(Response{ID: ids}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}