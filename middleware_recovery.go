package omnis

// Recovery -> recovers from any panics and returns JSON error message with StatusInternalServerError(500)
// func Recovery() gin.HandlerFunc {

// 	return func(ctx *gin.Context) {

// 		defer func() {

// 			if err := recover(); err != nil {

// 				var (
// 					render IRenderService = RenderService(ctx)
// 					goerr  *errors.Error  = errors.Wrap(err, 3)
// 					log    zerolog.Logger = defaultLogger().WithRequestContext(ctx).GetLogger().Level(zerolog.DebugLevel)
// 				)

// 				log.Err(goerr).Msg("")

// 				fmt.Print(string(goerr.Stack()))

// 				render.AsError(http.StatusInternalServerError, err)

// 				ctx.Abort()

// 				return
// 			}

// 		}()

// 		ctx.Next()

// 	}

// }
