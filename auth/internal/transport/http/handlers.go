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
	Login(ctx context.Context, login, password string) (string, error)
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

		token, err := srv.Login(r.Context(), req.Login, req.Password)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, 
		&http.Cookie{
			Name: "admin_token",
			Value: token,
			Path: "/",
			Domain: "",
			MaxAge: 86400,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
	})
}

