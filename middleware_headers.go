package omnis

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetHeaders(config *ServiceConfig) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// ctx.Header("content-type", "application/json")

		name := "omnis-service"
		version := "1.0.0"
		
		if config != nil {
			if config.Name != "" {
				name = config.Name
			}
			if config.Version != "" {
				version = config.Version
			}
		}

		ctx.Header("x-t3b-app", fmt.Sprintf("app:%s", name))
		ctx.Header("x-t3b-version", fmt.Sprintf("version:%s", version))

		ctx.Next()
	}

}
