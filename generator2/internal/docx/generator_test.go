package docx

import (
	"generator/internal/entity"
	"io"
	"os"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	file, _ := os.Open("../temp/ex.docx")
	_, _ = io.ReadAll(file)
	file.Close()

	type args struct {
		table entity.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGenerator_Generate",
			args: args{
				table: entity.Table{
					Columns: []string{"Column1", "Column2", "Column3"},
					Rows: [][]string{{"Row1", "Row2", "Row3"}, {"Row4", "Row5", "Row6"}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGenerator(Config{
				Key: "fa18c30cae8d384fdf028fe25f7ef998699694ffdfdca6ccd7a1508c335a330b",
			})
			if err != nil {
				t.Errorf("NewGenerator() error = %v", err)
			}
			
		})
	}
}
