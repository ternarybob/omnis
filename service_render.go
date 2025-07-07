package omnis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ternarybob/arbor"
	"github.com/ternarybob/funktion"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/phuslu/log"
)

type renderservice struct {
	ctx            *gin.Context
	internalLogger log.Logger
}

func RenderService(ctx *gin.Context) IRenderService {

	if ctx == nil {
		panic(fmt.Errorf("Context is nil"))
	}

	logger := defaultLogger()
	logger.Context = log.NewContext(nil).Str("function", "RenderService").Value()

	return &renderservice{
		ctx:            ctx,
		internalLogger: logger,
	}

}

func (s renderservice) AsResult(code int, payload interface{}) {

	output := s.getApiResponse(code)

	output.Result = payload

	s.respondwithJSON(code, output)

}

func (s renderservice) AsModel(code int, output interface{}) {

	apiresponse := s.getApiResponse(code)

	// Combine Api and Input Payloads
	apidata, err := json.Marshal(apiresponse)
	if err != nil {
		s.internalLogger.Warn().Msgf("Json Marshal err:%s", err.Error())
		return
	}

	if err := json.Unmarshal(apidata, &output); err != nil {
		s.internalLogger.Warn().Msgf("Json Marshal err:%s", err.Error())
		return
	}

	s.respondwithJSON(code, output)
}

func (s renderservice) AsResultWithError(code int, payload interface{}, err error) {

	output := s.getApiResponse(code)

	output.Result = payload

	if err != nil && cfg.Service.Scope == "DEV" {

		goerr := errors.Wrap(err, 3)

		output.Err = goerr.Error()
		output.Stack = funktion.SplitLines(string(goerr.Stack()))

	}

	s.respondwithJSON(code, output)

}

func (s renderservice) AsError(code int, err interface{}) {

	output := s.getApiResponse(code)

	if err != nil && cfg.Service.Scope == "DEV" {

		goerr := errors.Wrap(err, 3)

		output.Err = goerr.Error()
		output.Stack = funktion.SplitLines(string(goerr.Stack()))

	}

	s.respondwithJSON(code, output)

}

func (s renderservice) respondwithJSON(code int, payload interface{}) {

	if s.ctx == nil {
		panic(fmt.Errorf("Context is nil"))
	}

	s.ctx.Header("Content-Type", "application/json")

	if strings.ToUpper(cfg.Service.Scope) == "DEV" {
		s.ctx.IndentedJSON(code, payload)
		return

	}

	s.ctx.JSON(code, payload)

}

func (s renderservice) getApiResponse(code int) *ApiResponse {

	var (
		logs   = make(map[string]string)
		output = make(map[string]string)
	)

	if s.ctx == nil {
		panic(fmt.Errorf("Context is nil"))
	}

	s.internalLogger.Context = log.NewContext(nil).Str("function", "getApiResponse").Value()

	cid := s.getCorrelationID()

	if len(strings.TrimSpace(cid)) > 0 {
		// Get the default logger and retrieve memory logs
		logger := arbor.GetLogger()
		retrievedLogs, err := logger.GetMemoryLogs(cid, arbor.DebugLevel)
		if err != nil {
			logs["000"] = fmt.Sprintf("WRN|error retrieving logs %s", err)
		} else {
			logs = retrievedLogs
		}
	} else {
		// No correlation ID - add warning
		logs["000"] = "WRN|No correlation ID found - memory logging unavailable"
	}

	// Add "no logs found" warning if no logs are present
	if len(logs) == 0 {
		logs["000"] = fmt.Sprintf("WRN|No logs found for this request (memory logging may not be properly configured) CorrelationID:%s", cid)
	}

	if cfg.Service.Scope != "PRD" {
		output["url"] = s.ctx.FullPath()

		// Param
		for _, value := range s.ctx.Params {
			output[value.Key] = value.Value
		}

		// Form
		for key, value := range s.ctx.Request.PostForm {
			output[key] = strings.Join(value, ",")
		}

		// QueryString
		for key, value := range s.ctx.Request.URL.Query() {
			output[key] = strings.Join(value, ",")
		}
	}

	return &ApiResponse{
		Version:       cfg.Service.Version,
		Support:       cfg.Service.Support,
		Name:          cfg.Service.Name,
		Scope:         cfg.Service.Scope,
		Request:       output,
		Status:        code,
		CorrelationId: cid,
		Log:           logs,
	}

}

func (s renderservice) getCorrelationID() string {
	return GetCorrelationID(s.ctx)
}
