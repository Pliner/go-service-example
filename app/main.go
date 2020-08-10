package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DBNAME"),
		),
	}), &gorm.Config{
		Logger: NewGormLogger(logger),
	})

	if err != nil {
		logger.Fatal("Failed to open db", zap.Error(err))
		os.Exit(1)
	}

	usersRepository := &UsersRepository{db: db}
	usersController := &UsersController{usersRepository: usersRepository, logger: logger}

	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()

	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.MetricsPath = "/api/metrics"
	prometheus.Use(ginRouter)

	ginRouter.Use(RequestId())
	ginRouter.Use(Logging(logger))
	ginRouter.Use(Recovery(logger))
	ginRouter.GET("/api/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	ginRouter.GET("/api/users/:id", usersController.GetUserById)
	ginRouter.GET("/api/users", usersController.SelectUsers)
	ginRouter.POST("/api/users", usersController.SaveUser)
	ginRouter.DELETE("/api/users/:id", usersController.DeleteUserById)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: ginRouter,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Listen and served failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server was shut down")
}
