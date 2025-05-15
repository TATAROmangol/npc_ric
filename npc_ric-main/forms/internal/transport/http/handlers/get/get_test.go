package get

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"forms/internal/entities"
	"forms/internal/transport/http/handlers/get/mocks"
	"forms/pkg/logger"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetHandler_GetInstitutions(t *testing.T) {
	getter := mocks.NewMockGetter(gomock.NewController(t))

	type MockBehavior func()

	tests := []struct {
		name string
		MockBehavior
		wantResp   GetInstitutionsResponse
		wantStatus int
	}{
		{
			name: "valid institutions",
			MockBehavior: func() {
				getter.EXPECT().GetInstitutions(gomock.Any()).Return([]entities.Institution{
					{Id: 1, Name: "Institution 1", INN: 1234567890, Columns: []string{"Column1", "Column2"}},
					{Id: 2, Name: "Institution 2", INN: 9876543210, Columns: []string{"Column3", "Column4"}},
				}, nil)
			},
			wantResp: []entities.Institution{
				{Id: 1, Name: "Institution 1", INN: 1234567890, Columns: []string{"Column1", "Column2"}},
				{Id: 2, Name: "Institution 2", INN: 9876543210, Columns: []string{"Column3", "Column4"}},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error getting institutions",
			MockBehavior: func() {
				getter.EXPECT().GetInstitutions(gomock.Any()).Return(nil, errors.New("error getting institutions"))
			},
			wantResp:   nil,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior()

			h := NewHandler(getter)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/user/get", nil)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.GetInstitutions().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("GetInstitutions() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp GetInstitutionsResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("GetInstitutions() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("GetInstitutions() got = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetHandler_GetMentors(t *testing.T) {
	getter := mocks.NewMockGetter(gomock.NewController(t))

	type MockBehavior func()

	tests := []struct {
		name string
		MockBehavior
		wantResp   GetMentorsResponse
		wantStatus int
	}{
		{
			name: "valid mentors",
			MockBehavior: func() {
				getter.EXPECT().GetMentors(gomock.Any()).Return([]entities.Mentor{
					{Id: 1, Name: "Mentor 1"},
					{Id: 2, Name: "Mentor 2"},
				}, nil)
			},
			wantResp: []entities.Mentor{
				{Id: 1, Name: "Mentor 1"},
				{Id: 2, Name: "Mentor 2"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error getting mentors",
			MockBehavior: func() {
				getter.EXPECT().GetMentors(gomock.Any()).Return(nil, errors.New("error getting mentors"))
			},
			wantResp:   nil,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior()

			h := NewHandler(getter)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/user/get", nil)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.GetMentors().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("GetMentors() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp GetMentorsResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("GetMentors() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("GetMentors() got = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}

func TestGetHandler_GetInstitutionFromINN(t *testing.T) {
	getter := mocks.NewMockGetter(gomock.NewController(t))

	type MockBehavior func(inn int)

	tests := []struct {
		name string
		MockBehavior
		req        GetInstitutionFromINNRequest
		want       GetInstitutionFromINNResponse
		wantStatus int
	}{
		{
			name: "valid institution",
			MockBehavior: func(inn int) {
				getter.EXPECT().GetInstitutionFromINN(gomock.Any(), inn).Return(entities.Institution{
					Id:      1,
					Name:    "Test Institution",
					INN:     1234567890,
					Columns: []string{"Column1", "Column2"},
				}, nil)
			},
			req: GetInstitutionFromINNRequest{
				Inn: 1234567890,
			},
			want: GetInstitutionFromINNResponse{
				Id:      1,
				Name:    "Test Institution",
				INN:     1234567890,
				Columns: []string{"Column1", "Column2"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:         "invalid inn",
			MockBehavior: func(inn int) {},
			req: GetInstitutionFromINNRequest{
				Inn: 0,
			},
			want:       GetInstitutionFromINNResponse{},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req.Inn)

			h := NewHandler(getter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("GET", "/user/get", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.GetInstitutionFromINN().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("GetMentors() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp GetInstitutionFromINNResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("GetMentors() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("GetHandler.GetInstitutionFromINN() = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestGetHandler_GetFormColumns(t *testing.T) {
	getter := mocks.NewMockGetter(gomock.NewController(t))

	type MockBehavior func(institutionId int)

	tests := []struct {
		name   string
		MockBehavior
		req    GetFormColumnsRequest
		want   GetFormColumnsResponse
		wantStatus int
	}{
		{
			name: "valid columns",
			MockBehavior: func(institutionId int) {
				getter.EXPECT().GetFormColumns(gomock.Any(), institutionId).Return([]string{"Column1", "Column2"}, nil)
			},
			req: GetFormColumnsRequest{
				InstitutionId: 1,
			},
			want: []string{"Column1", "Column2"},
			wantStatus: http.StatusOK,
		},
		{
			name: "error getting columns",
			MockBehavior: func(institutionId int) {
				getter.EXPECT().GetFormColumns(gomock.Any(), institutionId).Return(nil, errors.New("error getting columns"))
			},
			req: GetFormColumnsRequest{
				InstitutionId: 1,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid institution id",
			MockBehavior: func(institutionId int) {},
			req: GetFormColumnsRequest{
				InstitutionId: 0,
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req.InstitutionId)

			h := NewHandler(getter)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("GET", "/user/get", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.GetFormColumns().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("GetFormColumns() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var resp GetFormColumnsResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("GetFormColumns() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("GetFormColumns() = %v, want %v", resp, tt.want)
			}
		})
	}
}
