package httpcreate

import (
	"context"
	"encoding/json"
	"net/http"
)

type Creater interface{
	GetInstitutionsId(ctx context.Context) ([]int, error)
}

type FromsHandler struct{
	Creater Creater
}

func (fh *FromsHandler) GetInstitutionsId() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := fh.Creater.GetInstitutionsId(r.Context())
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