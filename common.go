// -----------------------------------------------------------------------
// Last Modified: Wednesday, 27th August 2025 8:41:32 am
// Modified By: Bob McAllan
// -----------------------------------------------------------------------

package omnis

import (
	"github.com/phuslu/log"
)

const (
	CORRELATION_ID_KEY string = "correlationid"
	DEFAULT_TIMEFORMAT string = "01-02 15:04:05.000"
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
