package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Login   string `json:"login"`
	Password string `json:"password"`
}

type Loginer interface {
	Login(ctx context.Context, login, password string) (bool, error)
}

func LoginHandler(srv Loginer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ok, err := srv.Login(r.Context(), req.Login, req.Password)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "invalid login or password", http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

