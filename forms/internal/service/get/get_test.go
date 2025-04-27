package get

import (
	"context"
	"errors"
	"forms/internal/entities"
	"forms/internal/service/get/mocks"
	"forms/pkg/logger"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetService_GetInstitutions(t *testing.T) {
	repo := mocks.NewMockGetRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func()

	tests := []struct {
		name    string
		MockBehavior   MockBehavior
		want    []entities.Institution
		wantErr bool
	}{
		{
			name: "valid institutions",
			MockBehavior: func() {
				repo.EXPECT().GetInstitutions(gomock.Any()).Return([]entities.Institution{
					{Id: 1, Name: "Institution 1", INN: 1234567890, Columns: []string{"Column1", "Column2"}},
					{Id: 2, Name: "Institution 2", INN: 9876543210, Columns: []string{"Column3", "Column4"}},
				}, nil)
			},
			want: []entities.Institution{
				{Id: 1, Name: "Institution 1", INN: 1234567890, Columns: []string{"Column1", "Column2"}},
				{Id: 2, Name: "Institution 2", INN: 9876543210, Columns: []string{"Column3", "Column4"}},
			},
			wantErr: false,
		},
		{
			name: "error getting institutions",
			MockBehavior: func() {
				repo.EXPECT().GetInstitutions(gomock.Any()).Return(nil, errors.New("error getting institutions"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior()

			gs := &GetService{
				GetRepo: repo,
			}
			got, err := gs.GetInstitutions(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService.GetInstitutions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService.GetInstitutions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService_GetMentors(t *testing.T) {
	repo := mocks.NewMockGetRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func()

	tests := []struct {
		name    string
		MockBehavior   MockBehavior
		want    []entities.Mentor
		wantErr bool
	}{
		{
			name: "valid mentors",
			MockBehavior: func() {
				repo.EXPECT().GetMentors(gomock.Any()).Return([]entities.Mentor{
					{Id: 1, Info: "Mentor 1"},
					{Id: 2, Info: "Mentor 2"},
				}, nil)
			},
			want: []entities.Mentor{
				{Id: 1, Info: "Mentor 1"},
				{Id: 2, Info: "Mentor 2"},
			},
			wantErr: false,
		},
		{
			name: "error getting mentors",
			MockBehavior: func() {
				repo.EXPECT().GetMentors(gomock.Any()).Return(nil, errors.New("error getting mentors"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior()

			gs := &GetService{
				GetRepo: repo,
			}
			got, err := gs.GetMentors(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService.GetMentors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService.GetMentors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService_GetInstitutionFromINN(t *testing.T) {
	repo := mocks.NewMockGetRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(inn int)

	tests := []struct {
		name    string
		MockBehavior   MockBehavior
		inn     int
		want    entities.Institution
		wantErr bool
	}{
		{
			name: "valid institution",
			MockBehavior: func(inn int) {
				repo.EXPECT().GetInstitutionFromINN(gomock.Any(), inn).Return(entities.Institution{
					Id:   1,
					Name: "Institution 1",
					INN:  inn,
					Columns: []string{"Column1", "Column2"},
				}, nil)
			},
			inn: 1234567890,
			want: entities.Institution{
				Id:   1,
				Name: "Institution 1",
				INN:  1234567890,
				Columns: []string{"Column1", "Column2"},
			},
			wantErr: false,
		},
		{
			name: "error getting institution",
			MockBehavior: func(inn int) {
				repo.EXPECT().GetInstitutionFromINN(gomock.Any(), inn).Return(entities.Institution{}, errors.New("error getting institution"))
			},
			inn:     0,
			want:    entities.Institution{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.inn)

			gs := &GetService{
				GetRepo: repo,
			}
			got, err := gs.GetInstitutionFromINN(ctx, tt.inn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService.GetInstitutionFromINN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService.GetInstitutionFromINN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService_GetFormColumns(t *testing.T) {
	repo := mocks.NewMockGetRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int)

	tests := []struct {
		name    string
		MockBehavior   MockBehavior
		id     int
		want    []string
		wantErr bool
	}{
		{
			name: "valid columns",
			MockBehavior: func(id int) {
				repo.EXPECT().GetFormColumns(gomock.Any(), id).Return([]string{"Column1", "Column2"}, nil)
			},
			id:     1,
			want:    []string{"Column1", "Column2"},
			wantErr: false,
		},
		{
			name: "error getting columns",
			MockBehavior: func(id int) {
				repo.EXPECT().GetFormColumns(gomock.Any(), id).Return(nil, errors.New("error getting columns"))
			},
			id:     0,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.id)

			gs := &GetService{
				GetRepo: repo,
			}
			got, err := gs.GetFormColumns(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService.GetFormColumns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService.GetFormColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetService_GetFormRows(t *testing.T) {
	repo := mocks.NewMockGetRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int)

	tests := []struct {
		name    string
		MockBehavior   MockBehavior
		id     int
		want    []string
		wantErr bool
	}{
		{
			name: "valid rows",
			MockBehavior: func(id int) {
				repo.EXPECT().GetFormRows(gomock.Any(), id).Return([]string{"Row1", "Row2"}, nil)
			},
			id:     1,
			want:    []string{"Row1", "Row2"},
			wantErr: false,
		},
		{
			name: "error getting rows",
			MockBehavior: func(id int) {
				repo.EXPECT().GetFormRows(gomock.Any(), id).Return(nil, errors.New("error getting rows"))
			},
			id:     0,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.id)

			gs := &GetService{
				GetRepo: repo,
			}
			got, err := gs.GetFormRows(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetService.GetFormRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetService.GetFormRows() = %v, want %v", got, tt.want)
			}
		})
	}
}