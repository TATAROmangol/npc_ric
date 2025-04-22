package httpserver

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Deleter interface {
	DeleteInstitution() http.Handler
	DeleteMentor() http.Handler
}

type Getter interface {
	GetInstitutions() http.Handler
	GetMentors() http.Handler
	GetInstitutionFromINN() http.Handler
	GetFormColumns() http.Handler
}

type Putter interface {
	PutInstitutionInfo() http.Handler
	PutInstitutionColumns() http.Handler
	PutMentorRequest() http.Handler
}

type Poster interface {
	PostInstitution() http.Handler
	PostMentor() http.Handler
	PostForm() http.Handler
}

type Handlers struct{
	Deleter Deleter
	Getter Getter
	Putter Putter
	Poster Poster
}

type Middlewares interface {
	InitLoggerContextMiddleware(ctx context.Context) func(h http.Handler) http.Handler
	InitJsonContentTypeMiddleware() func(h http.Handler) http.Handler
}

type SendServer struct{
	ctx context.Context
	server *http.Server
}

func NewServer(ctx context.Context, cfg Config, h Handlers, m Middlewares) *SendServer {
	mux := mux.NewRouter()

	admin := mux.PathPrefix("/admin").Subrouter()
	admin.Use(m.InitLoggerContextMiddleware(ctx))
	admin.Use(m.InitJsonContentTypeMiddleware())
	admin.Handle("/post/institution", h.Poster.PostInstitution()).Methods(http.MethodPost)
	admin.Handle("/post/mentor", h.Poster.PostMentor()).Methods(http.MethodPost)
	admin.Handle("/put/institution", h.Putter.PutInstitutionInfo()).Methods(http.MethodPut)
	admin.Handle("/put/institution/columns", h.Putter.PutInstitutionColumns()).Methods(http.MethodPut)
	admin.Handle("/put/mentor", h.Putter.PutMentorRequest()).Methods(http.MethodPut)
	admin.Handle("/delete/institution", h.Deleter.DeleteInstitution()).Methods(http.MethodDelete)
	admin.Handle("/delete/mentor", h.Deleter.DeleteMentor()).Methods(http.MethodDelete)
	admin.Handle("/get/institutions", h.Getter.GetInstitutions()).Methods(http.MethodGet)
	admin.Handle("/get/institution", h.Getter.GetInstitutionFromINN()).Methods(http.MethodGet)
	admin.Handle("/get/mentors", h.Getter.GetMentors()).Methods(http.MethodGet)

	user := admin.PathPrefix("/user").Subrouter()
	user.Handle("/get/institutions", h.Getter.GetInstitutions()).Methods(http.MethodGet)
	user.Handle("/get/mentors", h.Getter.GetMentors()).Methods(http.MethodGet)
	user.Handle("/get/institution", h.Getter.GetInstitutionFromINN()).Methods(http.MethodGet)
	user.Handle("/get/form/columns", h.Getter.GetFormColumns()).Methods(http.MethodGet)
	user.Handle("post/form", h.Poster.PostForm()).Methods(http.MethodPost)
	
	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return &SendServer{
		ctx: ctx,
		server: server,
	}
}

func (s *SendServer) Run() error {
	return s.server.ListenAndServe()
}

func (s *SendServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}