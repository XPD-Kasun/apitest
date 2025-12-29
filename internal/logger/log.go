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

func InitLogger(logLevel string) {
	fmt.Println("Application logging init")
	setGlobalLogLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	appLogger.logger = &logger
}

func setGlobalLogLevel(logLevel string) {
	var logLvl zerolog.Level

	switch logLevel {
	case "DEBUG":
		logLvl = zerolog.DebugLevel
	case "INFO":
		logLvl = zerolog.InfoLevel
	case "WARN":
		logLvl = zerolog.WarnLevel
	case "ERROR":
		logLvl = zerolog.ErrorLevel
	case "FATAL":
		logLvl = zerolog.FatalLevel
	case "PANIC":
		logLvl = zerolog.PanicLevel
	case "DISABLE":
		logLvl = zerolog.Disabled
	case "TRACE":
		logLvl = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(logLvl)
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

func Trace() *zerolog.Event {
	return appLogger.logger.Trace()
}

func Fatal() *zerolog.Event {
	return appLogger.logger.Fatal()
}

func Root() *zerolog.Logger {
	return appLogger.logger
}
