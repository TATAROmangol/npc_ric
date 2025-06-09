package post

import (
	"context"
	"errors"
	"forms/internal/service/post/mocks"
	"forms/pkg/logger"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPostService_PostInstitution(t *testing.T) {
	repo := mocks.NewMockPostRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New(os.Stdout)
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

			ps := &PostService{
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
	l := logger.New(os.Stdout)
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

			ps := &PostService{
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
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(info []string, institutionId int)

	type args struct {
		info          []string
		institutionId int
		mentorId int
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
			MockBehavior: func(info []string, institutionId int) {
				repo.EXPECT().PostForm(gomock.Any(), info, institutionId).Return(1, nil)
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
			MockBehavior: func(info []string, institutionId int) {
				repo.EXPECT().PostForm(gomock.Any(), info, institutionId).Return(0, errors.New("error posting form"))
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
			tt.MockBehavior(tt.args.info, tt.args.institutionId)

			ps := &PostService{
				PostRepo: repo,
			}
			got, err := ps.PostForm(ctx, tt.args.info, tt.args.institutionId)
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
