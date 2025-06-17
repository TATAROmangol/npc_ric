package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Midlewarer interface{
	InitLogger() mux.MiddlewareFunc
	CheckAuth() mux.MiddlewareFunc
}

type Handlerer interface{
	UploadTemplate() http.Handler
	DeleteTemplate() http.Handler
	GetAllTemplates() http.Handler
}

type Server struct {
	srv *http.Server
}

func New(cfg *Config, m Midlewarer, h Handlerer) *Server {
	mux := mux.NewRouter()
	mux.Use(m.InitLogger())
	mux.Use(m.CheckAuth())
	mux.Handle("/upload", h.UploadTemplate()).Methods(http.MethodPost)
	mux.Handle("/template/{institution_id}", h.DeleteTemplate()).Methods(http.MethodDelete)
	mux.Handle("/", h.GetAllTemplates()).Methods(http.MethodGet)
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
