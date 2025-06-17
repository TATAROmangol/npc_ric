package httpserver

import "fmt"

type Config struct {
	Port int `env:"FORMS_HTTP_PORT"`
	Host string `env:"FORMS_HTTP_HOST"`
	AuthCookieName string `env:"AUTH_COOKIE_NAME"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}