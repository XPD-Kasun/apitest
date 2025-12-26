package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type AppLogger struct {
	logger *zerolog.Logger
}

var appLogger AppLogger

func InitLogger() {
	fmt.Println("Application logging init")
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	appLogger.logger = &logger
}

func Info() *zerolog.Event {
	return appLogger.logger.Info()
}

func Error() *zerolog.Event {
	return appLogger.logger.Error()
}

func Warn() *zerolog.Event {
	return appLogger.logger.Warn()
}

func Debug() *zerolog.Event {
	return appLogger.logger.Debug()
}

func Root() *zerolog.Logger {
	return appLogger.logger
}
