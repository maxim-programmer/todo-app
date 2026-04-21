// Package core_logger предоставляет структурированный логгер на базе go.uber.org/zap.
// Логи пишутся одновременно в stdout и в файл.
//
// Логгер передаётся через context.Context (паттерн «logger in context»),
// что позволяет автоматически добавлять к каждому сообщению
// request_id и другие поля, установленные в middleware.
package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loggerContextKey — приватный тип ключа для context.WithValue.
// Использование отдельного типа (а не string) исключает коллизии ключей
// с другими пакетами, которые тоже хранят данные в контексте.
type loggerContextKey struct{}

var (
	key = loggerContextKey{}
)

// Logger — обёртка над *zap.Logger, которая дополнительно хранит
// файловый дескриптор для корректного закрытия при завершении приложения.
type Logger struct {
	*zap.Logger

	file *os.File
}

// ToContext кладёт логгер в контекст. Вызывается в middleware Logger,
// чтобы все последующие обработчики могли получить логгер с request_id.
func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

// FromContext извлекает логгер из контекста.
// Паникует, если логгер не был добавлен — это программная ошибка,
// означающая, что middleware Logger не был подключён.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger создаёт логгер, который пишет в stdout и в файл одновременно.
// Каждый запуск создаёт новый лог-файл с именем вида "2006-01-02T15-04-05.000000.log".
//
// zapcore.NewTee объединяет несколько «ядер» (outputs) в одно:
// запись в одно ядро — автоматически запись во все.
func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	// ConsoleEncoder — человекочитаемый формат (не JSON), удобен для разработки.
	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// NewTee направляет логи одновременно в stdout и в файл.
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	// zap.AddCaller() добавляет к каждому сообщению имя файла и строку.
	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With создаёт дочерний логгер с дополнительными полями.
// Переопределяем метод, чтобы возвращать *core_logger.Logger (с файлом),
// а не базовый *zap.Logger.
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// Close закрывает файл логов. Должен вызываться через defer в main().
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
