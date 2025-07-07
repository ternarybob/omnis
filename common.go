package omnis

import (
	"github.com/phuslu/log"
	"github.com/ternarybob/satus"
)

const (
	CORRELATION_ID_KEY string = "correlationid"
	DEFAULT_TIMEFORMAT        = "01-02 15:04:05.000"
)

var (
	cfg *satus.AppConfig = satus.GetAppConfig()
)

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
