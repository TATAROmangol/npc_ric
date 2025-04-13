package httpcreate

import "fmt"

type Config struct {
	Port int `env:"FORMS_CREATE_HTTP_PORT"`
	Host string `env:"FORMS_CREATE_HTTP_HOST"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}