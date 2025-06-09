package post

import (
	"context"
	"os"
	"testing"

	"forms/internal/storage/repository/testcontainer"
	"forms/pkg/logger"
)

func TestPost_PostInstitution(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPost(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		args    struct {
			name    string
			inn     int
			columns []string
		}
		wantErr bool
	}{
		{
			name: "successful institution creation",
			args: struct {
				name    string
				inn     int
				columns []string
			}{
				name:    "Test Institution",
				inn:     444,
				columns: []string{"test1", "test2"},
			},
			wantErr: false,
		},
		{
			name: "duplicate INN",
			args: struct {
				name    string
				inn     int
				columns []string
			}{
				name:    "Another Institution",
				inn:     444,
				columns: []string{"test1", "test2"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.PostInstitution(ctx, tt.args.name, tt.args.inn, tt.args.columns)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.PostInstitution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPost_PostMentor(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPost(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	tests := []struct {
		name    string
		args    struct {
			name    string
		}
		wantErr bool
	}{
		{
			name: "successful mentor creation",
			args: struct {
				name string
			}{
				name: "Test Mentor",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.PostMentor(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.PostMentor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPost_PostForm(t *testing.T) {
	db, cleanup, err := testcontainer.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	defer cleanup()

	repo := NewPost(db)
	ctx := context.Background()
	l := logger.New(os.Stdout)
	ctx = logger.InitFromCtx(ctx, l)

	// Сначала создаем institution для теста формы
	institutionID, err := repo.PostInstitution(ctx, "Test Inst", 1234567890, []string{})
	if err != nil {
		t.Fatalf("failed to create test institution: %v", err)
	}

	tests := []struct {
		name          string
		institutionId int
		mentor_id	 int
		info          []string
		wantErr       bool
	}{
		{
			name:          "successful form creation",
			institutionId: institutionID,
			info:         []string{"field1", "field2"},
			wantErr:      false,
		},
		{
			name:          "invalid institution id",
			institutionId: -1,
			info:         []string{"field1"},
			wantErr:      true,
		},
		{
			name:          "nonexistent institution",
			institutionId: 999999,
			info:         []string{"field1"},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.PostForm(ctx, tt.info, tt.institutionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post.PostForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got <= 0 {
				t.Errorf("Post.PostForm() returned invalid ID: %v", got)
			}
		})
	}
}