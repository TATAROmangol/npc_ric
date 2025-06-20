package jwt

type Config struct {
	StringKey string `env: "JWT_KEY"`

}

func (c *Config) GetKey() []byte {
	return []byte(c.StringKey)
} 