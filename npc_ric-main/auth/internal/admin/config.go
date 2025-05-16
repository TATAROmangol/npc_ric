package admin

type Config struct {
	Login string `env:"ADMIN_LOGIN"`
	Password string `env:"ADMIN_PASSWORD"`
}