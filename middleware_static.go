package omnis

import (
	"time"

	"github.com/ternarybob/funktion"
	"github.com/ternarybob/satus"

	"github.com/gin-gonic/gin"
)

var (
	defaultExceptions []string = []string{"static/", "fav.ico", "favicon.ico", ".ico"}
)

func StaticRequests(cfg *satus.AppConfig, e []string) gin.HandlerFunc {

	log := warnLogger()

	return func(ctx *gin.Context) {

		requestExceptions := mergeUnique(defaultExceptions, e)

		if ctx.FullPath() != "" && funktion.ArrayContains(requestExceptions, ctx.FullPath()) {

			log.Trace().Msgf("Static Content")
			log.Trace().Msgf("path:%s contains:%t", ctx.FullPath(), funktion.ArrayContains(requestExceptions, ctx.FullPath()))

			if cfg.Service.Scope == "DEV" {
				ctx.Header("Expires", time.Now().Add(time.Minute*-1).Format(time.RFC3339))
			}

		}

		ctx.Next()

	}

}

func mergeUnique(arr1 []string, arr2 []string) []string {

	combined := append(arr1, arr2...)

	unique := make(map[string]bool)

	output := []string{}
	for _, val := range combined {
		if !unique[val] {
			unique[val] = true
			output = append(output, val)
		}
	}

	return output
}
