package service

import (
	"context"
	"errors"
	"forms/internal/service"
	"forms/internal/service/tests/mocks"
	"forms/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPostService_PostInstitution(t *testing.T) {
	repo := mocks.NewMockPostRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int)

	type args struct {
		name    string
		inn     int
		columns []string
	}

	tests := []struct {
		name    string
		MockBehavior  MockBehavior
		args    args
		id int
		want    int
		wantErr bool
	}{
		{
			name: "valid institution",
			MockBehavior: func(id int) {
				repo.EXPECT().PostInstitution(gomock.Any(), "Test Institution", 1234567890, []string{"Column1", "Column2"}).Return(id, nil)
			},
			args: args{
				name:    "Test Institution",
				inn:     1234567890,
				columns: []string{"Column1", "Column2"},
			},
			id: 1,
			want:    1,
			wantErr: false,
		},
		{
			name: "error posting institution",
			MockBehavior: func(id int) {
				repo.EXPECT().PostInstitution(gomock.Any(), "Test Institution", 1234567890, []string{"Column1", "Column2"}).Return(0, errors.New("error posting institution"))
			},
			args: args{
				name:    "Test Institution",
				inn:     1234567890,
				columns: []string{"Column1", "Column2"},
			},
			id: 0,
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.id)

			ps := &service.PostService{
				PostRepo: repo,
			}
			got, err := ps.PostInstitution(ctx, tt.args.name, tt.args.inn, tt.args.columns)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostService.PostInstitution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostService.PostInstitution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostService_PostMentor(t *testing.T) {
	repo := mocks.NewMockPostRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(name string)

	tests := []struct {
		name    string
		MockBehavior  MockBehavior
		nameMentor string
		id int
		want    int
		wantErr bool
	}{
		{
			name: "valid mentor",
			MockBehavior: func(name string) {
				repo.EXPECT().PostMentor(gomock.Any(), "Test Mentor").Return(1, nil)
			},
			nameMentor: "Test Mentor",
			id: 1,
			want:    1,
			wantErr: false,
		},
		{
			name: "error posting mentor",
			MockBehavior: func(name string) {
				repo.EXPECT().PostMentor(gomock.Any(), "Test Mentor").Return(0, errors.New("error posting mentor"))
			},
			nameMentor: "Test Mentor",
			id: 0,
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.nameMentor)

			ps := &service.PostService{
				PostRepo: repo,
			}
			got, err := ps.PostMentor(ctx, tt.nameMentor)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostService.PostMentor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostService.PostMentor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostService_PostForm(t *testing.T) {
	repo := mocks.NewMockPostRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int, info []string)

	type args struct {
		institutionId int
		info          []string
	}

	tests := []struct {
		name    string
		MockBehavior  MockBehavior
		args   args
		want    int
		wantErr bool
	}{
		{
			name: "valid form",
			MockBehavior: func(id int, info []string) {
				repo.EXPECT().PostForm(gomock.Any(), id, info).Return(1, nil)
			},
			args: args{
				institutionId: 1,
				info:          []string{"Column1", "Column2"},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error posting form",
			MockBehavior: func(id int, info []string) {
				repo.EXPECT().PostForm(gomock.Any(), id, info).Return(0, errors.New("error posting form"))
			},
			args: args{
				institutionId: 0,
				info:          []string{"Column1", "Column2"},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.args.institutionId, tt.args.info)

			ps := &service.PostService{
				PostRepo: repo,
			}
			got, err := ps.PostForm(ctx, tt.args.institutionId, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostService.PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PostService.PostForm() = %v, want %v", got, tt.want)
			}
		})
	}
}
