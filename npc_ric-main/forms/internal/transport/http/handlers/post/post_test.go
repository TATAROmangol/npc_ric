package post

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"forms/internal/transport/http/handlers/post/mocks"
	"forms/pkg/logger"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPostHandler_PostInstitution(t *testing.T) {
	poster := mocks.NewMockPoster(gomock.NewController(t))

	type MockBehavior func(req PostInstitutionRequest)

	tests := []struct {
		name string
		MockBehavior
		req        PostInstitutionRequest
		want       PostInstitutionResponse
		wantStatus int
	}{
		{
			name: "valid institution",
			MockBehavior: func(req PostInstitutionRequest) {
				poster.EXPECT().PostInstitution(gomock.Any(), req.Name, req.INN, req.Columns).Return(1, nil)
			},
			req: PostInstitutionRequest{
				Name:    "Test Institution",
				INN:     1234567890,
				Columns: []string{"Column1", "Column2"},
			},
			want: PostInstitutionResponse{
				Id: 1,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "error posting institution",
			MockBehavior: func(req PostInstitutionRequest) {
				poster.EXPECT().PostInstitution(gomock.Any(), req.Name, req.INN, req.Columns).Return(0, errors.New("error posting institution"))
			},
			req: PostInstitutionRequest{
				Name:    "Test Institution",
				INN:     1234567890,
				Columns: []string{"Column1", "Column2"},
			},
			want: PostInstitutionResponse{
				Id: 0,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:         "invalid name",
			MockBehavior: func(req PostInstitutionRequest) {},
			req: PostInstitutionRequest{
				Name:    "",
				INN:     1234567890,
				Columns: []string{"Column1", "Column2"},
			},
			want: PostInstitutionResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid inn",
			MockBehavior: func(req PostInstitutionRequest) {},
			req: PostInstitutionRequest{
				Name:    "Test Institution",
				INN:     0,
				Columns: []string{"Column1", "Column2"},
			},
			want: PostInstitutionResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "invalid columns",
			MockBehavior: func(req PostInstitutionRequest) {},
			req: PostInstitutionRequest{
				Name:    "Test Institution",
				INN:     1234567890,
				Columns: []string{},
			},
			want: PostInstitutionResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(poster)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PostInstitution().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("GetInstitutions() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var resp PostInstitutionResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("GetInstitutions() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("GetInstitutions() got = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestPostHandler_PostMentor(t *testing.T) {
	poster := mocks.NewMockPoster(gomock.NewController(t))

	type MockBehavior func(req PostMentorRequest)

	tests := []struct {
		name string
		MockBehavior
		req        PostMentorRequest
		want       PostMentorResponse
		wantStatus int
	}{
		{
			name: "valid mentor",
			MockBehavior: func(req PostMentorRequest) {
				poster.EXPECT().PostMentor(gomock.Any(), req.Name).Return(1, nil)
			},
			req: PostMentorRequest{
				Name: "Test",
			},
			want: PostMentorResponse{
				Id: 1,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "error posting mentor",
			MockBehavior: func(req PostMentorRequest) {
				poster.EXPECT().PostMentor(gomock.Any(), req.Name).Return(0, errors.New("error posting mentor"))
			},
			req: PostMentorRequest{
				Name: "test",
			},
			want: PostMentorResponse{
				Id: 0,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:         "invalid name",
			MockBehavior: func(req PostMentorRequest) {},
			req: PostMentorRequest{
				Name: "",
			},
			want: PostMentorResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(poster)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PostMentor().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("PostMentor() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var resp PostMentorResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("PostMentor() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("PostMentor() got = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestPostHandler_PostForm(t *testing.T) {
	poster := mocks.NewMockPoster(gomock.NewController(t))

	type MockBehavior func(req PostFormRequest)

	tests := []struct {
		name   string
		MockBehavior
		req PostFormRequest
		want   PostFormResponse
		wantStatus int
	}{
		{
			name: "valid form",
			MockBehavior: func(req PostFormRequest) {
				poster.EXPECT().GetFormColumns(gomock.Any(), req.InstitutionId).Return([]string{"a","b","c"}, nil)
				poster.EXPECT().PostForm(gomock.Any(), req.Info, req.InstitutionId).Return(1, nil)
			},
			req: PostFormRequest{
				Info: []string{"a", "b", "c"},
				InstitutionId: 1,
			},
			want: PostFormResponse{
				Id: 1,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid InstitutionId",
			MockBehavior: func(req PostFormRequest) {},
			req: PostFormRequest{
				Info: []string{"a", "b", "c"},
				InstitutionId: 0,
			},
			want: PostFormResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "err columns",
			MockBehavior: func(req PostFormRequest) {
				poster.EXPECT().GetFormColumns(gomock.Any(), req.InstitutionId).Return(nil, errors.New("test"))
			},
			req: PostFormRequest{
				Info: []string{"a", "b", "c"},
				InstitutionId: 1,
			},
			want: PostFormResponse{
				Id: 0,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "err len columns",
			MockBehavior: func(req PostFormRequest) {
				poster.EXPECT().GetFormColumns(gomock.Any(), req.InstitutionId).Return([]string{"a"}, nil)
			},
			req: PostFormRequest{
				Info: []string{"a", "b", "c"},
				InstitutionId: 1,
			},
			want: PostFormResponse{
				Id: 0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "err post form",
			MockBehavior: func(req PostFormRequest) {
				poster.EXPECT().GetFormColumns(gomock.Any(), req.InstitutionId).Return([]string{"a","b","c"}, nil)
				poster.EXPECT().PostForm(gomock.Any(), req.Info, req.InstitutionId).Return(0, errors.New("test"))
			},
			req: PostFormRequest{
				Info: []string{"a", "b", "c"},
				InstitutionId: 1,
			},
			want: PostFormResponse{
				Id: 0,
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.req)

			h := NewHandler(poster)

			rr := httptest.NewRecorder()

			body, _ := json.Marshal(tt.req)
			bytesBody := bytes.NewBuffer(body)
			req := httptest.NewRequest("POST", "/user/post", bytesBody)
			l := logger.New()
			ctx := logger.InitFromCtx(context.Background(), l)
			req = req.WithContext(ctx)

			h.PostForm().ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("PostForm() = %v, want %v", rr.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var resp PostFormResponse
			err := json.NewDecoder(rr.Result().Body).Decode(&resp)
			if err != nil {
				t.Errorf("PostForm() error decoding response = %v", err)
			}
			defer rr.Result().Body.Close()

			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("PostForm() got = %v, want %v", resp, tt.want)
			}
		})
	}
}
