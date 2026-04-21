// Package core_http_server содержит HTTP-сервер, систему маршрутизации
// и поддержку версионирования API.
package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config — конфигурация HTTP-сервера.
// Переменные окружения с префиксом "HTTP_": HTTP_ADDR, HTTP_SHUTDOWN_TIMEOUT, HTTP_ALLOWED_ORIGINS.
type Config struct {
	// Addr — адрес и порт сервера, например ":5050" (слушать на всех интерфейсах).
	Addr string `envconfig:"ADDR" required:"true"`

	// ShutdownTimeout — время ожидания завершения активных запросов при остановке сервера.
	// По истечении этого времени сервер принудительно закрывает соединения.
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`

	// AllowedOrigins — список разрешённых origins для CORS.
	// envconfig парсит строку с запятой как слайс: "http://a.com,http://b.com".
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" required:"true"`
}

// NewConfig читает конфигурацию сервера из переменных окружения.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust — «Must»-вариант конструктора: паникует при ошибке.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}

	return config
}
