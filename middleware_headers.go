package omnis

import (
	"fmt"

	"github.com/ternarybob/satus"

	"github.com/gin-gonic/gin"
)

func SetHeaders(cfg *satus.AppConfig) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// ctx.Header("content-type", "application/json")

		ctx.Header("x-t3b-app", fmt.Sprintf("app:%s", cfg.Service.Name))
		ctx.Header("x-t3b-version", fmt.Sprintf("version:%s", cfg.Service.Version))

		ctx.Next()
	}

}
