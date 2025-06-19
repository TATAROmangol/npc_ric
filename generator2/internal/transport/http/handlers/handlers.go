package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"generator/pkg/logger"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//go:generate mockgen -source=handlers.go -destination=./mocks/mock_handlers.go -package=mocks

type Srv interface {
	DeleteTemplate(ctx context.Context, id int) error
	UploadTemplate(ctx context.Context, id int, file multipart.File) error
	GenerateTemplate(ctx context.Context, id int) (*os.File, func(), error)
}

type Handlers struct{
	srv Srv
}

func New(srv Srv) *Handlers {
	return &Handlers{
		srv: srv,
	}
}

func (h *Handlers) DeleteTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid id", err)
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := h.srv.DeleteTemplate(r.Context(), id); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to delete template", err)
			http.Error(w, "failed to delete template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

func (h *Handlers) UploadTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil{
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to parse form", err)
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		institutionIdStr := r.FormValue("institution_id")
		if institutionIdStr == "" {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid institution_id", nil)
			http.Error(w, "invalid institution_id", http.StatusBadRequest)
			return
		}
		institutionId, err := strconv.Atoi(institutionIdStr)
		if err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid institution_id", err)
			http.Error(w, "invalid institution_id", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "File upload error", err)
			http.Error(w, "File upload error: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(header.Filename, ".docx") {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "unknown format file", nil)
			http.Error(w, "Only .docx files are allowed", http.StatusBadRequest)
			return
		}

		if err := h.srv.UploadTemplate(r.Context(), institutionId, file); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed to upload file", err)
			http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusAccepted)
	})
}

func (h *Handlers) GenerateTemplate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := struct{
			InstitutionId int `json:"institution_id"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "invalid response", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid request"))
			return
		}

		doc, remove, err := h.srv.GenerateTemplate(r.Context(), req.InstitutionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer doc.Close()
		defer remove()

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", doc.Name()))
    	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

		if _, err := doc.Stat(); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed get stat file", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed get stat file" + err.Error()))
			return
		}
		doc.Seek(0, 0)

		if _, err := io.Copy(w, bufio.NewReader(doc)); err != nil {
			logger.GetFromCtx(r.Context()).ErrorContext(r.Context(), "failed return for download file", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
}
