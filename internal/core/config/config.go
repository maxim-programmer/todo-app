// Package core_config содержит общую конфигурацию приложения,
// не привязанную к конкретной инфраструктурной компоненте.
package core_config

import (
	"fmt"
	"os"
	"time"
)

// Config хранит общие настройки приложения.
type Config struct {
	// TimeZone задаёт часовой пояс для time.Local во всём приложении.
	// Все time.Now() будут возвращать время в этом часовом поясе.
	TimeZone *time.Location
}

// NewConfig читает конфигурацию из переменных окружения.
// Переменная TIME_ZONE принимает IANA-идентификатор часового пояса:
//   - UTC
//   - Europe/Berlin
//   - America/New_York
//   - Europe/Moscow
func NewConfig() (*Config, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("load time zone: %s: %w", tz, err)
	}

	return &Config{
		TimeZone: zone,
	}, nil
}

// NewConfigMust — «Must»-вариант конструктора: паникует при ошибке.
// Используется при старте приложения, когда работа без конфигурации невозможна.
func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
