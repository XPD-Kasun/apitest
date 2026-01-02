package logger

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
)

func TestSuit(t *testing.T) {

	InitLogger("DEBUG")

	t.Run("TestIsAppLoggerReturnedFromCtx", func(t *testing.T) {

		var r = zerolog.Ctx(context.Background())
		if r != appLogger.logger {
			t.Error("zerolog.Ctx does not return the appLogger for an empty context")
		}
		if r == nil {
			t.Error("zerolog.Ctx returns empty logger")
		}
	})

	t.Run("TestCtx", func(t *testing.T) {
		appLogger.logger.Info().Msg("sdf")
	})

}
