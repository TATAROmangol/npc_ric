package put

import (
	"context"
	"errors"
	"forms/internal/service/put/mocks"
	"forms/pkg/logger"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPutService_PutInstitutionInfo(t *testing.T) {
	repo := mocks.NewMockPutRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int, name string, inn int)

	type args struct {
		id   int
		name string
		inn  int
	}
	tests := []struct {
		name    string
		MockBehavior MockBehavior
		args    args
		wantErr bool
	}{
		{
			name: "valid institution info",
			MockBehavior: func(id int, name string, inn int) {
				repo.EXPECT().PutInstitutionInfo(gomock.Any(), id, name, inn).Return(nil)
			},
			args: args{
				id:   1,
				name: "Test Institution",
				inn:  1234567890,
			},
			wantErr: false,
		},
		{
			name: "error putting institution info",
			MockBehavior: func(id int, name string, inn int) {
				repo.EXPECT().PutInstitutionInfo(gomock.Any(), id, name, inn).Return(errors.New("error putting institution info"))
			},
			args: args{
				id:   1,
				name: "Test Institution",
				inn:  1234567890,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.args.id, tt.args.name, tt.args.inn)

			ps := &PutService{
				PutRepo: repo,
			}
			if err := ps.PutInstitutionInfo(ctx, tt.args.id, tt.args.name, tt.args.inn); (err != nil) != tt.wantErr {
				t.Errorf("PutService.PutInstitutionInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPutService_PutInstitutionColumns(t *testing.T) {
	repo := mocks.NewMockPutRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int, columns []string)

	type args struct {
		id   int
		columns  []string
	}
	tests := []struct {
		name    string
		MockBehavior MockBehavior
		args    args
		wantErr bool
	}{
		{
			name: "valid institution columns",
			MockBehavior: func(id int, columns []string) {
				repo.EXPECT().DeleteForms(gomock.Any(), id).Return(nil)
				repo.EXPECT().PutInstitutionColumns(gomock.Any(), id, columns).Return(nil)
			},
			args: args{
				id:   1,
				columns:  []string{"Column1", "Column2"},
			},
			wantErr: false,
		},
		{
			name: "error putting institution columns",
			MockBehavior: func(id int, columns []string) {
				repo.EXPECT().DeleteForms(gomock.Any(), id).Return(nil)
				repo.EXPECT().PutInstitutionColumns(gomock.Any(), id, columns).Return(errors.New("error putting institution columns"))
			},
			args: args{
				id:   1,
				columns:  []string{"Column1", "Column2"},
			},
			wantErr: true,
		},
		{
			name: "error deleting forms",
			MockBehavior: func(id int, columns []string) {
				repo.EXPECT().DeleteForms(gomock.Any(), id).Return(errors.New("error deleting forms"))
			},
			args: args{
				id:   1,
				columns:  []string{"Column1", "Column2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.args.id, tt.args.columns)

			ps := &PutService{
				PutRepo: repo,
			}
			if err := ps.PutInstitutionColumns(ctx, tt.args.id, tt.args.columns); (err != nil) != tt.wantErr {
				t.Errorf("PutService.PutInstitutionColumns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPutService_PutMentor(t *testing.T) {
	repo := mocks.NewMockPutRepo(gomock.NewController(t))

	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	type MockBehavior func(id int, info string)

	type args struct {
		id   int
		info  string
	}
	tests := []struct {
		name    string
		MockBehavior MockBehavior
		args    args
		wantErr bool
	}{
		{
			name: "valid mentor info",
			MockBehavior: func(id int, info string) {
				repo.EXPECT().PutMentor(gomock.Any(), id, info).Return(nil)
			},
			args: args{
				id:   1,
				info: "Test Mentor Info",
			},
			wantErr: false,
		},
		{
			name: "error putting mentor info",
			MockBehavior: func(id int, info string) {
				repo.EXPECT().PutMentor(gomock.Any(), id, info).Return(errors.New("error putting mentor info"))
			},
			args: args{
				id:   1,
				info: "Test Mentor Info",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.MockBehavior(tt.args.id, tt.args.info)

			ps := &PutService{
				PutRepo: repo,
			}
			if err := ps.PutMentor(ctx, tt.args.id, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("PutService.PutMentor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
