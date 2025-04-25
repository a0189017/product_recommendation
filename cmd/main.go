package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"product_recommendation/migrations"
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/routine"
	"product_recommendation/pkg/server"
	"product_recommendation/pkg/server/middleware"
	"product_recommendation/pkg/server/routers"

	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

func main() {
	_ = godotenv.Load()

	config := config.GetConfig()

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel = signal.NotifyContext(context.Background(), constants.SignalsToShutdown...)
	defer cancel()

	dbConfig := config.DB
	db := database.New(&database.DBOptions{
		User:          dbConfig.User,
		Password:      dbConfig.Password,
		Host:          dbConfig.Host,
		Port:          dbConfig.Port,
		Name:          dbConfig.Name,
		SlowThreshold: dbConfig.SlowThreshold,
		Colorful:      dbConfig.Log.Colorful,
	})

	redis := database.NewRedisClient()
	redisClient := redis.GetRedisClient()

	// run migrations
	err := migrations.RunMigration(db.GetDB())
	if err != nil {
		log.Fatalf("%v", err)
	}

	db.LoadSchemaFields()

	// cronjob
	go routine.ProductRecommendation(db, redisClient)

	c := cron.New(cron.WithSeconds())
	initCronJob(c, config, db, redisClient)
	c.Start()

	r := server.New(db, redisClient)

	versionGroup := r.Group("/v" + constants.CurrentVersion[0:1])
	// auth
	authGroup := versionGroup.Group("/auth")
	authGroup.POST("/register", routers.Register)
	authGroup.POST("/login", routers.Login)
	// TODO: otp can be resend
	authGroup.POST("/otp/verify", middleware.VerifyToken(constants.LoginTypeOtp), routers.VerifyOTP)

	// product
	productGroup := versionGroup.Group("/product")
	productGroup.Use(middleware.VerifyToken(constants.LoginTypeToken))
	productGroup.GET("/recommendation", routers.GetProductRecommendation)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server start failed: %s\n", err)
		}
	}()

	<-ctx.Done()

	cancel()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// gracefully shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func initCronJob(c *cron.Cron, config *config.Config, db database.DBInterface, redis *redis.Client) {
	_, err := c.AddFunc(config.CronJobSpec.SyncRecommendation, func() { routine.ProductRecommendation(db, redis) })
	if err != nil {
		log.Fatalf("add corn job for check cron is work failed, err = %v", err)
	}
}
