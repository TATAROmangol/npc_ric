package admin

import "testing"

func TestAdmin_IsValid(t *testing.T) {
	cfg := Config{
		Login:    "admin",
		Password: "password",
	}
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{
			name: "valid credentials",
			args: args{
				login:    "admin",
				password: "password",
			},
			want: true,
		},
		{
			name: "wrong login",
			args: args{
				login:    "user",
				password: "password",
			},
			want: false,
		},
		{
			name: "empty credentials",
			args: args{
				login:    "",
				password: "",
			},
			want: false,
		},
		{
			name: "wrong password",
			args: args{
				login:    "admin",
				password: "wrongpassword",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Admin{
				cfg: cfg,
			}
			if got := a.IsValid(tt.args.login, tt.args.password); got != tt.want {
				t.Errorf("Admin.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
