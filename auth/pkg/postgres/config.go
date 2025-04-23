package postgres

import (
	"fmt"
)

type Config struct {
	Host     string `env:"PG_HOST"`
	Port     string `env:"PG_PORT"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASSWORD"`
	DBName   string `env:"PG_DB_NAME"`
	SSL      string `env:"PG_SSL"`
}

func (c *Config) GetConnString() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSL,
	)
}