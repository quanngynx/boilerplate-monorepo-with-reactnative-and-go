package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/cache"
	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewHealthHandler(db *gorm.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redisClient}
}

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
	Redis  string `json:"redis"`
}

func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := database.Ping(h.db); err != nil {
		dbStatus = "error"
	}

	redisStatus := "ok"
	if err := cache.Ping(ctx, h.redis); err != nil {
		redisStatus = "error"
	}

	status := http.StatusOK
	overall := "ok"
	if dbStatus != "ok" || redisStatus != "ok" {
		status = http.StatusServiceUnavailable
		overall = "degraded"
	}

	c.JSON(status, healthResponse{
		Status: overall,
		DB:     dbStatus,
		Redis:  redisStatus,
	})
}
