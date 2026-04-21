package core_goredis_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config — конфигурация Redis-клиента.
// Все поля читаются из переменных окружения с префиксом REDIS_.
type Config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     string        `envconfig:"PORT" default:"6379"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	DB       int           `envconfig:"DB" default:"0"`
	TTL      time.Duration `envconfig:"TTL" default:"5m"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get Redis pool config: %w", err))
	}

	return config
}
