package httpserver

import "fmt"

type Config struct {
	Port int `env:"FORMS_HTTP_PORT"`
	Host string `env:"FORMS_HTTP_HOST"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}