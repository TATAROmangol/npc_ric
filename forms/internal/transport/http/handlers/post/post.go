package post

import (
	"context"
	"encoding/json"
	"fmt"
	"forms/pkg/logger"
	"net/http"
)

var (
	ErrInvalidInn		 = fmt.Errorf("invalid Inn")
	ErrInvalidName		 = fmt.Errorf("invalid Name")
	ErrInvalidColumns	 = fmt.Errorf("invalid Columns")
	ErrInvalidInstitutionId = fmt.Errorf("invalid InstitutionId")
)

//go:generate mockgen -source=post.go -destination=./mocks/mock_post.go -package=mocks

type Post interface {
	PostInstitution(ctx context.Context, name string, inn int, columns []string) (int, error)
	PostMentor(ctx context.Context, name string) (int, error)
	PostForm(ctx context.Context, info []string, institutionId int) (int, error)
}

type Get interface{
	GetFormColumns(ctx context.Context, id int) ([]string, error)
}

type Poster interface{
	Post
	Get
}

type PostHandler struct {
	srv Poster
}

func NewHandler(srv Poster) *PostHandler {
	return &PostHandler{
		srv: srv,
	}
}

func (h *PostHandler) PostInstitution() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PostInstitutionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostInstitution handler", ErrInvalidName)
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		if req.INN <= 0 {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostInstitution handler", ErrInvalidInn)
			http.Error(w, "inn is required", http.StatusBadRequest)
			return
		}

		if len(req.Columns) == 0 {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostInstitution handler", ErrInvalidColumns)
			http.Error(w, "columns are required", http.StatusBadRequest)
			return
		}

		id, err := h.srv.PostInstitution(r.Context(), req.Name, req.INN, req.Columns)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(PostInstitutionResponse{Id: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostInstitution handler", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *PostHandler) PostMentor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PostMentorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostMentor handler", ErrInvalidName)
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		id, err := h.srv.PostMentor(r.Context(), req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(PostMentorResponse{Id: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostMentor handler", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func (h *PostHandler) PostForm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PostFormRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed service in PostForm handler", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.InstitutionId <= 0 {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostForm handler", ErrInvalidInstitutionId)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		columns, err := h.srv.GetFormColumns(r.Context(), req.InstitutionId)
		if err != nil{
			http.Error(w, "failed service in PostForm", http.StatusInternalServerError)
			return
		}

		if len(columns) != len(req.Info){
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		id, err := h.srv.PostForm(r.Context(), req.Info, req.InstitutionId)
		if err != nil {
			http.Error(w, "invalid request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(PostFormResponse{Id: id}); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed in PostForm handler", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
