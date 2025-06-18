package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Midlewarer interface{
	InitLoggerCtx() mux.MiddlewareFunc
	Operation() mux.MiddlewareFunc
	InitJsonContentType() mux.MiddlewareFunc
	CheckAuth(cookieName string) mux.MiddlewareFunc
}

type Handlerer interface{
	UploadTemplate() http.Handler
	DeleteTemplate() http.Handler
	GenerateTemplate() http.Handler
}

type Server struct {
	srv *http.Server
}

func New(cfg *Config, m Midlewarer, h Handlerer) *Server {
	mux := mux.NewRouter()
	mux.Use(m.InitLoggerCtx())
	mux.Use(m.Operation())
	mux.Use(m.InitJsonContentType())
	mux.Use(m.CheckAuth(cfg.AuthCookieName))
	mux.Handle("/upload", h.UploadTemplate()).Methods(http.MethodPost)
	mux.Handle("/template/{institution_id}", h.DeleteTemplate()).Methods(http.MethodDelete)
	mux.Handle("/generate", h.GenerateTemplate()).Methods(http.MethodPost)

	return &Server{
		srv: &http.Server{
			Addr:    cfg.Host,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Close()
}
