package omnis

import "github.com/ternarybob/arbor"

type IRenderService interface {
	AsModel(code int, output interface{})
	AsResult(code int, payload interface{})
	AsResultWithError(code int, payload interface{}, err error)
	AsError(code int, err interface{})
	WithLogger(logger arbor.ILogger) IRenderService
	WithConfig(config *ServiceConfig) IRenderService
}
