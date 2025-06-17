package verify

import "fmt"

type Config struct {
	Port int `env:"GENERATOR_AUTH_GRPC_PORT"`
	Host string `env:"GENERATOR_AUTH_GRPC_HOST"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}