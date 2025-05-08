package put

import (
	"context"
	"forms/internal/storage/repository/get"
	"forms/internal/storage/repository/testcontainer"
	"forms/pkg/logger"
	"testing"
)

func TestPut_PutInstitutionInfo(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPut(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	type args struct {
		name string
		inn  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful institution update",
			args: args{
				name: "Updated Institution",
				inn:  1234567890,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.PutInstitutionInfo(ctx, 1, tt.args.name, tt.args.inn); (err != nil) != tt.wantErr {
				t.Errorf("Put.PutInstitutionInfo() error = %v, wantErr %v", err, tt.wantErr)
			}

			getRepo := get.NewGet(db)
			institution, err := getRepo.GetInstitutionFromINN(ctx, tt.args.inn)
			if err != nil {
				t.Errorf("Get.GetInstitutionFromINN() error = %v", err)
				return
			}
			if institution.Name != tt.args.name {
				t.Errorf("Get.GetInstitutionFromINN() name got = %v, want %v", institution.Name, tt.args.name)
			}
		})
	}
}

func TestPut_PutInstitutionColumns(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPut(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		columns []string
		wantErr bool
	}{
		{
			name:    "successful institution columns update",
			columns: []string{"Column1", "Column2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.PutInstitutionColumns(ctx, 1, tt.columns); (err != nil) != tt.wantErr {
				t.Errorf("Put.PutInstitutionColumns() error = %v, wantErr %v", err, tt.wantErr)
			}

			getRepo := get.NewGet(db)
			columns, err := getRepo.GetFormColumns(ctx, 1)
			if err != nil {
				t.Errorf("Get.GetFormColumns() error = %v", err)
				return
			}
			if len(columns) != len(tt.columns) {
				t.Errorf("Get.GetFormColumns() columns length got = %v, want %v", len(columns), len(tt.columns))
			}
			for i, column := range columns {
				if column != tt.columns[i] {
					t.Errorf("Get.GetFormColumns() column got = %v, want %v", column, tt.columns[i])
				}
			}
		})
	}
}

func TestPut_PutMentor(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPut(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	t.Run("first", func(t *testing.T) {
		if err := repo.PutMentor(ctx, 1, "first"); err != nil {
			t.Errorf("Put.PutMentor() error = %v", err)
		}

		getRepo := get.NewGet(db)
		mentors, err := getRepo.GetMentors(ctx)
		if err != nil {
			t.Errorf("Get.PutMentor() error = %v", err)
			return
		}
		for _, m := range mentors {
			if m.Id == 1 && m.Name != "first" {
				t.Errorf("Get.PutMentor() id got = %v, want %v", m.Id, 1)
			}
		}
	})
}