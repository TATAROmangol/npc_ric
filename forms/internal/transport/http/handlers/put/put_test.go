package put

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"forms/internal/transport/http/handlers/put/mocks"
	"forms/pkg/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPutHandler_PutInstitutionInfo(t *testing.T) {
	putter := mocks.NewMockPutter(gomock.NewController(t))

	type MockBehavior func(req PutInstitutionInfoRequest)

	tests := []struct {
		name string
		MockBehavior
		req  PutInstitutionInfoRequest
		want int
	}{
		{
			name: "valid institution info",
			MockBehavior: func(req PutInstitutionInfoRequest) {
				putter.EXPECT().PutInstitutionInfo(gomock.Any(), req.Id, req.Name, req.INN).Return(nil)
			},
			req: PutInstitutionInfoRequest{
				Id:   1,
				Name: "Test Institution",
				INN:  1234567890,
			},
			want: http.StatusOK,
		},
		{
			name: "error updating institution info",
			MockBehavior: func(req PutInstitutionInfoRequest) {
				putter.EXPECT().PutInstitutionInfo(gomock.Any(), req.Id, req.Name, req.INN).Return(errors.New("error updating institution info"))
			},
			req: PutInstitutionInfoRequest{
				Id:   1,
				Name: "Test Institution",
				INN:  1234567890,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(putter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New(os.Stdout)
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PutInstitutionInfo().ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("GetInstitutions() = %v, want %v", rr.Code, tt.want)
			}
		})
	}
}

func TestPutHandler_PutInstitutionColumns(t *testing.T) {
	putter := mocks.NewMockPutter(gomock.NewController(t))

	type MockBehavior func(req PutInstitutionColumnsRequest)

	tests := []struct {
		name   string
		MockBehavior
		req   PutInstitutionColumnsRequest
		want   int
	}{
		{
			name: "valid institution columns",
			MockBehavior: func(req PutInstitutionColumnsRequest) {
				putter.EXPECT().PutInstitutionColumns(gomock.Any(), req.Id, req.Columns).Return(nil)
			},
			req: PutInstitutionColumnsRequest{
				Id:      1,
				Columns: []string{"column1", "column2"},
			},
			want: http.StatusOK,
		},
		{
			name: "error updating institution columns",
			MockBehavior: func(req PutInstitutionColumnsRequest) {
				putter.EXPECT().PutInstitutionColumns(gomock.Any(), req.Id, req.Columns).Return(errors.New("error updating institution columns"))
			},
			req: PutInstitutionColumnsRequest{
				Id:      1,
				Columns: []string{"column1", "column2"},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(putter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New(os.Stdout)
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PutInstitutionColumns().ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("GetInstitutions() = %v, want %v", rr.Code, tt.want)
			}
		})
	}
}

func TestPutHandler_PutMentor(t *testing.T) {
	putter := mocks.NewMockPutter(gomock.NewController(t))

	type MockBehavior func(req PutMentorRequest)

	tests := []struct {
		name   string
		MockBehavior
		req    PutMentorRequest
		want   int
	}{
		{
			name: "valid mentor",
			MockBehavior: func(req PutMentorRequest) {
				putter.EXPECT().PutMentor(gomock.Any(), req.Id, req.Name).Return(nil)
			},
			req: PutMentorRequest{
				Id:   1,
				Name: "Test Mentor",
			},
			want: http.StatusOK,
		},
		{
			name: "error updating mentor",
			MockBehavior: func(req PutMentorRequest) {
				putter.EXPECT().PutMentor(gomock.Any(), req.Id, req.Name).Return(errors.New("error updating mentor"))
			},
			req: PutMentorRequest{
				Id:   1,
				Name: "Test Mentor",
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(putter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New(os.Stdout)
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PutMentor().ServeHTTP(rr, req)

			if rr.Code != tt.want {
				t.Errorf("GetInstitutions() = %v, want %v", rr.Code, tt.want)
			}
		})
	}
}
