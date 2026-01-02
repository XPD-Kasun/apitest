package logger

import (
	"apitest/internal/core/common"
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var appLogger AppLogger

type AppLogger struct {
	logger *zerolog.Logger
}

type logHook struct {
}

// Run implements [zerolog.Hook].
func (l logHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	r := e.GetCtx()
	if r == nil {
		return
	}
	if r.Value(common.AppContextKey) == nil {
		return
	}

	appReqCtx := r.Value(common.AppContextKey).(*common.AppRequestContext)
	if appReqCtx.CorrelationId != "" {
		e.Str("corrid", string(appReqCtx.CorrelationId))
	}
}

func InitLogger(logLevel string) {
	fmt.Println("Application logging init")
	setGlobalLogLevel(logLevel)
	logger := zerolog.New(os.Stdout).Hook(logHook{}).With().Timestamp().Logger()

	appLogger.logger = &logger
	zerolog.DefaultContextLogger = appLogger.logger
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

func SetCorrelationID(ctx context.Context, cid common.Uniqueid) {
	appReqCtx := ctx.Value(common.AppContextKey).(*common.AppRequestContext)
	appReqCtx.CorrelationId = cid
}
