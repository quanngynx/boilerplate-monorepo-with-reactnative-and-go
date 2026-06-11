package middleware

import (
	"log/slog"

	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/apperror"
	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/logger"
	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(log *slog.Logger, debugMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		appErr := apperror.FromError(err)

		switch appErr.Code {
		case "VALIDATION_FAILED":
			logger.InputValidationFail(log, RequestID(c), appErr.Message)
		}

		if appErr.Err != nil {
			log.Error("request failed",
				slog.String("request_id", RequestID(c)),
				slog.String("code", appErr.Code),
				slog.Any("error", appErr.Err),
			)
		}

		response.WriteProblem(c, appErr, debugMode)
	}
}

func Recovery(log *slog.Logger, debugMode bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				var cause error
				switch value := recovered.(type) {
				case error:
					cause = value
				case string:
					cause = apperror.Internal(value, nil)
				default:
					cause = apperror.Internal("unexpected panic", nil)
				}

				logger.SysCrash(log, cause)
				appErr := apperror.Internal("An error occurred, please retry", cause)
				response.WriteProblem(c, appErr, debugMode)
				c.Abort()
			}
		}()

		c.Next()
	}
}
