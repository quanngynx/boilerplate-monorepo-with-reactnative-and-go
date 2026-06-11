package response

import (
	"net/http"

	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/apperror"
	"github.com/gin-gonic/gin"
)

type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
	Code     string `json:"code,omitempty"`
}

func WriteProblem(c *gin.Context, appErr *apperror.AppError, includeDebug bool) {
	problem := ProblemDetail{
		Type:     "about:blank",
		Title:    http.StatusText(appErr.HTTPStatus),
		Status:   appErr.HTTPStatus,
		Detail:   appErr.Message,
		Instance: c.Request.URL.Path,
		Code:     appErr.Code,
	}

	if includeDebug && appErr.Err != nil {
		problem.Detail = appErr.Err.Error()
	}

	c.Header("X-ERROR", "true")
	c.Header("Content-Type", "application/problem+json")
	c.JSON(appErr.HTTPStatus, problem)
}
