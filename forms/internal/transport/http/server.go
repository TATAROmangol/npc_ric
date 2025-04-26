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
	PutMentor() http.Handler
}

type Poster interface {
	PostInstitution() http.Handler
	PostMentor() http.Handler
	PostForm() http.Handler
}

type Handlers interface{
	Deleter
	Getter
	Putter
	Poster
}

type Verifier interface {
	Verify(ctx context.Context, token string) (bool, error)
}

type Middlewares interface {
	InitLoggerContextMiddleware(ctx context.Context) func(h http.Handler) http.Handler
	InitJsonContentTypeMiddleware() func(h http.Handler) http.Handler
	OperationMiddleware() func(h http.Handler) http.Handler
	AuthMiddleware() func(h http.Handler) http.Handler
}

type HTTPServer struct{
	ctx context.Context
	server *http.Server
}

func NewServer(ctx context.Context, cfg Config, h Handlers, m Middlewares) *HTTPServer {
	mux := mux.NewRouter()
	mux.Use(m.InitLoggerContextMiddleware(ctx))
	mux.Use(m.InitJsonContentTypeMiddleware())
	mux.Use(m.OperationMiddleware())

	admin := mux.PathPrefix("/admin").Subrouter()
	admin.Use(m.AuthMiddleware())
	admin.Handle("/post/institution", h.PostInstitution()).Methods(http.MethodPost)
	admin.Handle("/post/mentor", h.PostMentor()).Methods(http.MethodPost)
	admin.Handle("/put/institution", h.PutInstitutionInfo()).Methods(http.MethodPut)
	admin.Handle("/put/institution/columns", h.PutInstitutionColumns()).Methods(http.MethodPut)
	admin.Handle("/put/mentor", h.PutMentor()).Methods(http.MethodPut)
	admin.Handle("/delete/institution", h.DeleteInstitution()).Methods(http.MethodDelete)
	admin.Handle("/delete/mentor", h.DeleteMentor()).Methods(http.MethodDelete)
	admin.Handle("/get/institutions", h.GetInstitutions()).Methods(http.MethodGet)
	admin.Handle("/get/institution", h.GetInstitutionFromINN()).Methods(http.MethodGet)
	admin.Handle("/get/mentors", h.GetMentors()).Methods(http.MethodGet)

	user := admin.PathPrefix("/user").Subrouter()
	user.Handle("/get/institutions", h.GetInstitutions()).Methods(http.MethodGet)
	user.Handle("/get/mentors", h.GetMentors()).Methods(http.MethodGet)
	user.Handle("/get/institution", h.GetInstitutionFromINN()).Methods(http.MethodGet)
	user.Handle("/get/form/columns", h.GetFormColumns()).Methods(http.MethodGet)
	user.Handle("/post/form", h.PostForm()).Methods(http.MethodPost)
	
	server := &http.Server{
		Addr:    cfg.Addr(),
		Handler: mux,
	}

	return &HTTPServer{
		ctx: ctx,
		server: server,
	}
}

func (hs *HTTPServer) Run() error {
	return hs.server.ListenAndServe()
}

func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.server.Shutdown(ctx)
}