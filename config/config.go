package config

type Config struct {
	Secret string
}

func New(secret string) *Config {
	return &Config{Secret: secret}
}
