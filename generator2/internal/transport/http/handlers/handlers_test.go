package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"generator/internal/transport/http/handlers/mocks"
	"generator/pkg/logger"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_DeleteTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockSrv(ctrl)
	l := logger.New(os.Stdout)

	tests := []struct {
		name         string
		id           string
		mockBehavior func()
		wantStatus   int
		wantBody     string
	}{
		{
			name: "successful deletion",
			id:   "123",
			mockBehavior: func() {
				srv.EXPECT().DeleteTemplate(123).Return(nil)
			},
			wantStatus: http.StatusAccepted,
		},
		{
			name: "invalid id - not a number",
			id:   "abc",
			mockBehavior: func() {
				// No expectations for invalid ID
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid id",
		},
		{
			name: "service error",
			id:   "123",
			mockBehavior: func() {
				srv.EXPECT().DeleteTemplate(123).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed to delete template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBehavior != nil {
				tt.mockBehavior()
			}

			h := New(srv)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/templates/"+tt.id, nil)
			
			// Set context with logger
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)
			
			// Set URL vars for mux
			vars := map[string]string{
				"id": tt.id,
			}
			req = mux.SetURLVars(req, vars)

			h.DeleteTemplate().ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != "" {
				assert.Contains(t, rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestHandlers_UploadTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockSrv(ctrl)
	l := logger.New(os.Stdout)

	tests := []struct {
		name           string
		institutionID string
		fileContent   string
		fileName      string
		mockBehavior  func()
		wantStatus    int
		wantBody      string
	}{
		{
			name:          "successful upload",
			institutionID: "123",
			fileContent:   "test content",
			fileName:      "test.docx",
			mockBehavior: func() {
				srv.EXPECT().UploadTemplate(123, gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusAccepted,
		},
		{
			name:          "invalid institution_id",
			institutionID: "abc",
			mockBehavior:  func() {},
			wantStatus:    http.StatusBadRequest,
			wantBody:      "invalid institution_id",
		},
		{
			name:          "missing file",
			institutionID: "123",
			mockBehavior:  func() {},
			wantStatus:    http.StatusBadRequest,
			wantBody:      "File upload error",
		},
		{
			name:          "invalid file format",
			institutionID: "123",
			fileContent:   "test content",
			fileName:      "test.txt",
			mockBehavior:  func() {},
			wantStatus:    http.StatusBadRequest,
			wantBody:      "Only .docx files are allowed",
		},
		{
			name:          "service error",
			institutionID: "123",
			fileContent:   "test content",
			fileName:      "test.docx",
			mockBehavior: func() {
				srv.EXPECT().UploadTemplate(123, gomock.Any()).Return(errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Failed to upload file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBehavior != nil {
				tt.mockBehavior()
			}

			h := New(srv)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Добавляем institution_id
			_ = writer.WriteField("institution_id", tt.institutionID)

			// Добавляем файл, если указан
			if tt.fileContent != "" && tt.fileName != "" {
				part, _ := writer.CreateFormFile("file", tt.fileName)
				_, _ = part.Write([]byte(tt.fileContent))
			}

			writer.Close()

			req := httptest.NewRequest("POST", "/templates", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			
			// Set context with logger
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			h.UploadTemplate().ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != "" {
				assert.Contains(t, rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestHandlers_GenerateTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mocks.NewMockSrv(ctrl)
	l := logger.New(os.Stdout)

	tests := []struct {
		name         string
		requestBody  interface{}
		mockBehavior func()
		wantStatus   int
		wantBody     string
		wantHeaders  map[string]string
	}{
		{
			name: "successful generation",
			requestBody: map[string]int{
				"institution_id": 123,
			},
			mockBehavior: func() {
				file, _ := os.CreateTemp("", "test*.docx")
				_, _ = file.WriteString("test content")
				_, _ = file.Seek(0, 0)
				srv.EXPECT().GenerateTemplate(123).Return(file, nil)
			},
			wantStatus: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Disposition": "attachment; filename=",
				"Content-Type":        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
		},
		{
			name:        "invalid request",
			requestBody: "invalid",
			mockBehavior: func() {
				// No expectations for invalid request
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request",
		},
		{
			name: "service error",
			requestBody: map[string]int{
				"institution_id": 123,
			},
			mockBehavior: func() {
				srv.EXPECT().GenerateTemplate(123).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "service error",
		},
		{
			name: "file copy error",
			requestBody: map[string]int{
				"institution_id": 123,
			},
			mockBehavior: func() {
				file, _ := os.CreateTemp("", "test*.docx")
				file.Close() 
				srv.EXPECT().GenerateTemplate(123).Return(file, nil)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "failed get stat file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockBehavior != nil {
				tt.mockBehavior()
			}

			h := New(srv)

			var reqBody []byte
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("GET", "/templates/generate", bytes.NewBuffer(reqBody))
			
			// Set context with logger
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			h.GenerateTemplate().ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != "" {
				assert.Contains(t, rr.Body.String(), tt.wantBody)
			}

			// Проверяем заголовки для успешного случая
			if tt.wantHeaders != nil {
				for key, val := range tt.wantHeaders {
					assert.Contains(t, rr.Header().Get(key), val)
				}
			}
		})
	}
}