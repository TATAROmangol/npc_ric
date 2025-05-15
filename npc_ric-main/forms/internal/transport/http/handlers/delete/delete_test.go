package delete

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"forms/internal/service/delete/mocks"
	"forms/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestDeleteHandler_DeleteInstitution(t *testing.T) {
	deleter := mocks.NewMockDeleteRepo(gomock.NewController(t))

	type MockBehavior func(id int)

	tests := []struct {
		name         string
		mockBehavior MockBehavior
		request      DeleteInstitutionRequest
		wantStatus   int
	}{
		{
			name: "successful deletion",
			mockBehavior: func(id int) {
				deleter.EXPECT().DeleteInstitution(gomock.Any(), id).Return(nil)
			},
			request: DeleteInstitutionRequest{
				Id: 1,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error deletion",
			mockBehavior: func(id int) {
				deleter.EXPECT().DeleteInstitution(gomock.Any(), id).Return(errors.New("error deleting institution"))
			},
			request: DeleteInstitutionRequest{
				Id: 1,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid id",
			mockBehavior: func(id int) {},
			request: DeleteInstitutionRequest{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.request.Id)

			h := NewHandler(deleter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.request)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("DELETE", "/user/delete", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.DeleteInstitution().ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("DeleteHandler.DeleteInstitution() got = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteHandler_DeleteMentor(t *testing.T) {
	deleter := mocks.NewMockDeleteRepo(gomock.NewController(t))

	type MockBehavior func(id int)

	tests := []struct {
		name   string
		mockBehavior MockBehavior
		request      DeleteMentorRequest
		wantStatus   int
	}{
		{
			name: "successful deletion",
			mockBehavior: func(id int) {
				deleter.EXPECT().DeleteMentor(gomock.Any(), id).Return(nil)
			},
			request: DeleteMentorRequest{
				Id: 1,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error deletion",
			mockBehavior: func(id int) {
				deleter.EXPECT().DeleteMentor(gomock.Any(), id).Return(errors.New("error deleting mentor"))
			},
			request: DeleteMentorRequest{
				Id: 1,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid id",
			mockBehavior: func(id int) {},
			request: DeleteMentorRequest{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.request.Id)

			h := NewHandler(deleter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.request)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("DELETE", "/user/delete", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.DeleteMentor().ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("DeleteHandler.DeleteInstitution() got = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
