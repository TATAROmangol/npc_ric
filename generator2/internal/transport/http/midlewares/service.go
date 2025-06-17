package midlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (m *Midleware) InitLogger() mux.MiddlewareFunc{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Midleware) CheckAuth() mux.MiddlewareFunc{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}