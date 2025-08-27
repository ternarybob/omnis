// -----------------------------------------------------------------------
// Last Modified: Wednesday, 27th August 2025 8:49:17 am
// Modified By: Bob McAllan
// -----------------------------------------------------------------------

package omnis

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/ternarybob/arbor"
	"github.com/ternarybob/funktion"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/phuslu/log"
)

type renderservice struct {
	ctx            *gin.Context
	internalLogger log.Logger
	logger         arbor.ILogger
	config         *ServiceConfig
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

func (s *renderservice) WithLogger(logger arbor.ILogger) IRenderService {
	s.logger = logger
	return s
}

func (s *renderservice) WithConfig(config *ServiceConfig) IRenderService {
	s.config = config
	return s
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

	if err != nil && s.getScope() == "DEV" {

		goerr := errors.Wrap(err, 3)

		output.Err = goerr.Error()
		output.Stack = funktion.SplitLines(string(goerr.Stack()))

	}

	s.respondwithJSON(code, output)

}

func (s renderservice) AsError(code int, err interface{}) {

	output := s.getApiResponse(code)

	if err != nil && s.getScope() == "DEV" {

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

	if strings.ToUpper(s.getScope()) == "DEV" {
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
		// Use provided logger if available, otherwise fall back to default
		var loggerToUse arbor.ILogger
		if s.logger != nil {
			loggerToUse = s.logger
		} else {
			loggerToUse = arbor.GetLogger()
		}

		retrievedLogs, err := loggerToUse.GetMemoryLogs(cid, arbor.DebugLevel)
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

	if s.getScope() != "PRD" {
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
		Version:       s.getVersion(),
		Name:          s.getName(),
		Scope:         s.getScope(),
		Request:       output,
		Status:        code,
		CorrelationId: cid,
		Log:           logs,
	}

}

func (s renderservice) getCorrelationID() string {
	return GetCorrelationID(s.ctx)
}

// Configuration helper methods with defaults
func (s *renderservice) getVersion() string {
	baseVersion := "1.0.0"
	if s.config != nil && s.config.Version != "" {
		baseVersion = s.config.Version
	}

	// Add build information: version+build.goversion.timestamp
	buildTime := time.Now().Format("20060102.150405")
	goVersion := runtime.Version()

	return fmt.Sprintf("%s+build.%s.%s", baseVersion, goVersion, buildTime)
}

func (s *renderservice) getName() string {
	if s.config != nil && s.config.Name != "" {
		return s.config.Name
	}
	return "omnis-service"
}

func (s *renderservice) getScope() string {
	if s.config != nil && s.config.Scope != "" {
		return s.config.Scope
	}
	return "DEV"
}
