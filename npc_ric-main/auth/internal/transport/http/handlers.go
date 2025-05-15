package httpserver

import (
	"auth/pkg/logger"
	"context"
	"encoding/json"
	"net/http"
)

//go:generate mockgen -source=handlers.go -destination=./mocks/mock_handlers.go -package=mocks

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

		if req.Login == "" || req.Password == "" {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid login or password", nil)
			http.Error(w, "login and password are required", http.StatusBadRequest)
			return
		}

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

func LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_token")
		if err != nil || cookie == nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "not found cookie", nil)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusAccepted)
	})
}

