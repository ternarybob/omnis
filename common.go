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

func defaultLogger() arbor.ILogger { return arbor.Logger() }

func warnLogger() arbor.ILogger { return arbor.Logger().WithLevel(arbor.WarnLevel) }
