package table

import "fmt"

type Config struct {
	Port int `env:"GENERATOR_TABLE_GRPC_PORT"`
	Host string `env:"GENERATOR_TABLE_GRPC_HOST"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}