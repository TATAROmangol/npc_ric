package get

import (
	"context"
	"forms/internal/entities"
	"forms/internal/storage/repository/testcontainer"
	"forms/pkg/logger"
	"os"
	"testing"
)

func TestGet_GetInstitutions(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewGet(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		want    []entities.Institution
		wantErr bool
	}{
		{
			name: "successful institution retrieval",
			want: []entities.Institution{
				{Id: 1, Name: "A", INN: 111},
				{Id: 2, Name: "B", INN: 222},
				{Id: 3, Name: "C", INN: 333},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			institutions, err := repo.GetInstitutions(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.GetInstitutions() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i, institution := range institutions {
				if institution.Name != tt.want[i].Name {
					t.Errorf("Get.GetInstitutions() name got = %v, want %v", institution.Name, tt.want[i].Name)
				}
				if institution.INN != tt.want[i].INN {
					t.Errorf("Get.GetInstitutions() inn got = %v, want %v", institution.INN, tt.want[i].INN)
				}
			}
		})
	}
}

func TestGet_GetMentors(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewGet(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		want    []entities.Mentor
		wantErr bool
	}{
		{
			name: "successful mentor retrieval",
			want: []entities.Mentor{
				{Id: 1, Name: "A"},
				{Id: 2, Name: "B"},
				{Id: 3, Name: "C"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mentors, err := repo.GetMentors(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.GetMentors() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i, mentor := range mentors {
				if mentor.Name != tt.want[i].Name {
					t.Errorf("Get.GetInstitutions() name got = %v, want %v", mentor.Name, tt.want[i].Name)
				}
			}
		})
	}
}

func TestGet_GetInstitutionFromINN(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewGet(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		inn     int
		want    entities.Institution
		wantErr bool
	}{
		{
			name: "successful institution retrieval by INN",
			inn:  111,
			want: entities.Institution{
				Id:   1,
				Name: "A",
				INN:  111,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			institution, err := repo.GetInstitutionFromINN(ctx, tt.inn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.GetInstitutionFromINN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if institution.Name != tt.want.Name {
				t.Errorf("Get.GetInstitutionFromINN() name got = %v, want %v", institution.Name, tt.want.Name)
			}
		})
	}
}

func TestGet_GetFormColumns(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewGet(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		id      int
		want    []string
		wantErr bool
	}{
		{
			name:    "successful form columns retrieval",
			id:      1,
			want:    []string{"1", "2", "3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns, err := repo.GetFormColumns(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.GetFormColumns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i, column := range columns {
				if column != tt.want[i] {
					t.Errorf("Get.GetFormColumns() column got = %v, want %v", column, tt.want[i])
				}
			}
		})
	}
}

func TestGet_GetFormRows(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewGet(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		institution_id  int
		want    [][]string
		wantErr bool
	}{
		{
			name:           "successful form rows retrieval",
			institution_id: 1,
			want: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
				{"G", "H", "I"},

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetFormRows(ctx, tt.institution_id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get.GetFormRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i, row := range got {
				for j, value := range row {
					if value != tt.want[i][j] {
						t.Errorf("Get.GetFormRows() row[%d][%d] got = %v, want %v", i, j, value, tt.want[i][j])
					}
				}
			}
		})
	}
}
