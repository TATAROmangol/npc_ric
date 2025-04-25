package migrator

import (
	"fmt"
)

type Config struct{
	User     string `env:"MIGRATE_USER"`
	Password string `env:"MIGRATE_PASSWORD"`
	Host     string	 `env:"MIGRATE_HOST"`
	Port     string `env:"MIGRATE_PORT"`
	DBName   string `env:"MIGRATE_DB_NAME"`
	SSL      string `env:"MIGRATE_SSL"`
	Schema   string `env:"MIGRATE_SCHEMA"`
}

func (c *Config) GetConnString() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v&x-migrations-table=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSL, c. Schema,
	)
}

