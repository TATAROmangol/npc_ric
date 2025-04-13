package httpsend

import (
	"context"
	"encoding/json"
	"net/http"
)

type Sender interface{
	SendForm(ctx context.Context, institution string, info []string) (int, error)
}

type SendHandler struct{
	Sender Sender
}

func (fh *SendHandler) SendForm() http.Handler{
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

		id, err := fh.Sender.SendForm(r.Context(), req.Institution, req.Info)
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