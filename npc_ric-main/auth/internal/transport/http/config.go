package httpserver

import "fmt"

type Config struct {
	Port int `env:"AUTH_HTTP_PORT"`
	Host string `env:"AUTH_HTTP_HOST"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}