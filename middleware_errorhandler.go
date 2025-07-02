package omnis

// func ErrorHandler() echo.HandlerFunc {

// 	return func(ctx echo.Context) error {

// 		log := defaultLogger().WithRequestContext(ctx).GetLogger().Level(zerolog.DebugLevel)

// 		err := ctx.Errors.Last()

// 		if err == nil {
// 			ctx.Next()
// 			return nil
// 		}

// 		log.Debug().Msg("Errors Detected")

// 		render := RenderService(ctx)

// 		log.Warn().Err(err).Msg("")

// 		render.AsError(http.StatusInternalServerError, err)

// 		ctx.Abort()

// 		return nil

// 	}

// }
