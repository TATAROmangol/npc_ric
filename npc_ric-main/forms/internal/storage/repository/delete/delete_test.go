package delete

import (
	"context"
	"forms/internal/storage/repository/get"
	"forms/internal/storage/repository/testcontainer"
	"forms/pkg/logger"
	"testing"
)

func TestDelete_DeleteInstitution(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewDelete(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name          string
		institutionId int
		wantErr       bool
	}{
		{
			name:          "successful institution deletion",
			institutionId: 1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.DeleteInstitution(ctx, tt.institutionId); (err != nil) != tt.wantErr {
				t.Errorf("Delete.DeleteInstitution() error = %v, wantErr %v", err, tt.wantErr)
			}

			getRepo := get.NewGet(db)
			institutions, _ := getRepo.GetInstitutions(ctx)
			if len(institutions) != 2 {
				t.Errorf("Get.GetInstitutions() length got = %v, want %v", len(institutions), 2)
			}

			forms, err := getRepo.GetFormRows(ctx, tt.institutionId)
			if err != nil {
				t.Errorf("Get.GetFormRows() error = %v", err)
			}
			if len(forms) != 0 {
				t.Errorf("Get.GetFormRows() length got = %v, want %v", len(forms), 0)
			}
		})
	}
}

func TestDelete_DeleteMentor(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewDelete(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name     string
		mentorId int
		wantErr  bool
	}{
		{
			name:     "successful mentor deletion",
			mentorId: 1,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.DeleteMentor(ctx, tt.mentorId); (err != nil) != tt.wantErr {
				t.Errorf("Delete.DeleteMentor() error = %v, wantErr %v", err, tt.wantErr)
			}

			getRepo := get.NewGet(db)
			mentors, _ := getRepo.GetMentors(ctx)
			if len(mentors) != 2 {
				t.Errorf("Get.GetMentors() length got = %v, want %v", len(mentors), 2)
			}
		})
	}
}

func TestDelete_DeleteForms(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewDelete(db)
	ctx := context.Background()
	l := logger.New()
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		institution_id int
		wantErr bool
	}{
		{
			name:          "successful forms deletion",
			institution_id: 1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.DeleteForms(ctx, tt.institution_id); (err != nil) != tt.wantErr {
				t.Errorf("Delete.DeleteForms() error = %v, wantErr %v", err, tt.wantErr)
			}

			getRepo := get.NewGet(db)
			forms, err := getRepo.GetFormRows(ctx, tt.institution_id)
			if err != nil {
				t.Errorf("Get.GetFormRows() error = %v", err)
				return
			}
			if len(forms) != 0 {
				t.Errorf("Get.GetFormRows() length got = %v, want %v", len(forms), 0)
			}
		})
	}
}
