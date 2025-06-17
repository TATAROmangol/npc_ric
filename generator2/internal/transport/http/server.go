package http

import (
	"generator2/internal/transport/http/handlers"
	"generator2/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct{
	srv *http.Server
}

func New(cfg *Config, l *logger.Logger) *Server{
	mux := mux.NewRouter()
	mux.Handle("/upload", handlers.UploadFile())
	return &Server{
		srv: &http.Server{
			Addr: cfg.Host,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error{
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error{
	return s.srv.Close()
}
