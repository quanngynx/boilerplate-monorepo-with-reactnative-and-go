package router

import (
	"log/slog"

	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/config"
	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/handler"
	"github.com/GenAI-Fund/boilerplate-monorepo-with-reactnative-and-go/apps/server/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	Config *config.Config
	Log    *slog.Logger
	DB     *gorm.DB
	Redis  *redis.Client
}

func New(deps Dependencies) *gin.Engine {
	gin.SetMode(deps.Config.GinMode)

	debugMode := deps.Config.GinMode == gin.DebugMode

	engine := gin.New()
	engine.Use(middleware.Recovery(deps.Log, debugMode))
	engine.Use(middleware.RequestLogger(deps.Log))
	engine.Use(middleware.CORS(deps.Config.AllowedOrigins))
	engine.Use(middleware.ErrorHandler(deps.Log, debugMode))

	healthHandler := handler.NewHealthHandler(deps.DB, deps.Redis)
	helloHandler := handler.NewHelloHandler()

	engine.GET("/health", healthHandler.Health)

	v1 := engine.Group("/api/v1")
	{
		v1.GET("/hello", helloHandler.Hello)
	}

	return engine
}
