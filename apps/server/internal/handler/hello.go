package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HelloHandler struct{}

func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

type helloResponse struct {
	Message string `json:"message"`
}

func (h *HelloHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, helloResponse{
		Message: "Hello from Go Gin",
	})
}
