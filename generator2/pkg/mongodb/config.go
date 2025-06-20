package mongodb

import "fmt"

type Config struct {
	Host           string `env:"MONGO_HOST"`
	Port           string `env:"MONGO_PORT"`
	User           string `env:"MONGO_USER"`
	Password       string `env:"MONGO_PASSWORD"`
	DBName         string `env:"MONGO_DB_NAME"`
	CollectionName string `env:"MONGO_COLLECTION_NAME"`
	AuthSource     string `env:"MONGO_AUTH_SOURCE"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("mongodb://%v:%v@%v:%v/%v?authSource=%v", c.User, c.Password, c.Host, c.Port, c.DBName, c.AuthSource)
}
