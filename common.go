package omnis

import (
	"github.com/ternarybob/arbor"
	"github.com/ternarybob/satus"
)

const (
	CORRELATION_ID_KEY string = "correlationid"
)

var (
	cfg *satus.AppConfig = satus.GetAppConfig()
)

func defaultLogger() arbor.IConsoleLogger { return arbor.ConsoleLogger() }

func warnLogger() arbor.IConsoleLogger { return arbor.ConsoleLogger().WithLevel(arbor.WarnLevel) }
