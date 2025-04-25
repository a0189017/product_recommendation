package server

import (
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/server/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func New(database database.DBInterface, redis *redis.Client) *gin.Engine {
	c := config.GetConfig()
	gin.SetMode(gin.DebugMode)

	engine := gin.New()
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(middleware.SetLogger)
	if database != nil {
		engine.Use(middleware.SetDatabase(database))
	}
	if redis != nil {
		engine.Use(middleware.SetRedis(redis))
	}
	engine.Use(middleware.Recover())

	if c.Server.EnableCORS {
		config := cors.Config{
			AllowMethods: []string{
				"PUT",
				"PATCH",
				"OPTIONS",
				"POST",
				"GET",
				"DELETE",
			},
			AllowHeaders: []string{
				"Content-Type",
				"X-Amz-Date",
				"Authorization",
				"X-Api-Key",
				"X-Amz-Security-Token",
			},
			AllowOrigins:     c.Server.AllowOrigins,
			MaxAge:           15 * time.Minute,
			AllowCredentials: false,
		}
		engine.Use(cors.New(config))
	}

	return engine
}
