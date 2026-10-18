package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		HTTP    http
		PG      pg
		JWT     jwt
		Swagger swagger
	}
	// HTTP -.
	http struct {
		Port string `env:"HTTP_PORT,required"`
	}

	// PG -.
	pg struct {
		URL string `env:"PG_URL,required"`
	}

	// JWT -.
	jwt struct {
		Secret      string        `env:"JWT_SECRET,required"`
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}

	// Swagger -.
	swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	// 2. Теперь caarlos0/env считывает переменные из окружения в структуру cfg.
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
