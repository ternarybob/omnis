package omnis

import (
	"github.com/ternarybob/arbor"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type gincorrelation struct {
	ctx    gin.Context
	logger arbor.IConsoleLogger
}

func SetCorrelationID() gin.HandlerFunc {

	log := warnLogger().WithFunction().GetLogger()

	return func(ctx *gin.Context) {

		_, exists := ctx.Get(CORRELATION_ID_KEY)

		if !exists {

			uuid, err := uuid.NewRandom()
			if err != nil {
				log.Warn().Err(err).Msg("")
				return
			}

			ctx.Set(CORRELATION_ID_KEY, uuid.String())

		}

		// Test the get
		output := ctx.MustGet(CORRELATION_ID_KEY).(string)

		// Write to Response Header
		ctx.Header(CORRELATION_ID_KEY, output)

		log.Debug().Msgf("CorrelationId:%s", output)

	}

}

func GetCorrelationID(c *gin.Context) string {
	return c.GetString(CORRELATION_ID_KEY)
}
