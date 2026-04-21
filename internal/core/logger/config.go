package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config хранит настройки логгера, читаемые из переменных окружения.
// Библиотека kelseyhightower/envconfig автоматически заполняет поля из env-переменных
// с префиксом "LOGGER_": LOGGER_LEVEL, LOGGER_FOLDER.
type Config struct {
	// Level — минимальный уровень логирования: DEBUG, INFO, WARN, ERROR.
	Level string `envconfig:"LEVEL" default:"DEBUG"`

	// Folder — директория, в которой будут создаваться файлы логов.
	// Каждый запуск приложения создаёт новый файл с timestamp в имени.
	Folder string `envconfig:"FOLDER" required:"true"`
}

// NewConfig читает конфигурацию логгера из переменных окружения.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust — «Must»-вариант конструктора: паникует при ошибке.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}

	return config
}
