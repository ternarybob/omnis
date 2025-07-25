package omnis

import (
	"github.com/phuslu/log"
	"github.com/ternarybob/arbor"
	"github.com/ternarybob/arbor/models"
	"github.com/ternarybob/satus"
)

const (
	CORRELATION_ID_KEY string = "correlationid"
	DEFAULT_TIMEFORMAT        = "01-02 15:04:05.000"
)

var (
	cfg *satus.AppConfig = satus.GetAppConfig()
)

// getArborLogger returns a configured arbor logger using satus configuration
func getArborLogger() arbor.ILogger {
	return arbor.Logger().
		WithConsoleWriter(models.WriterConfiguration{
			Type: models.LogWriterTypeConsole,
		}).
		WithLevelFromString(satus.GetLogLevel()).
		WithPrefix("omnis")
}

func defaultLogger() log.Logger {
	return log.Logger{
		Level:      log.DebugLevel,
		TimeFormat: DEFAULT_TIMEFORMAT,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}
}

func warnLogger() log.Logger {
	return log.Logger{
		Level:      log.WarnLevel,
		TimeFormat: DEFAULT_TIMEFORMAT,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}
}
