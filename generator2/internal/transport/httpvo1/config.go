package httpvo1

import "fmt"

type Config struct {
	Host string `env:"GENERATOR_HTTP_HOST"`
	Port int    `env:"GENERATOR_HTTP_PORT"`
	AuthCookieName string `env:"AUTH_COOKIE_NAME"`
}

func (c *Config) Addr() string{
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}
