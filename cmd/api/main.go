// @title SmallSupermarket APIs
// @version 1.0
// @description SmallSupermarket APIs
// @host localhost:8080
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"supermarket-backend/infrastructure/db"
	"supermarket-backend/infrastructure/jwt"
	"supermarket-backend/internal/config"
	"supermarket-backend/internal/handler"
	"supermarket-backend/internal/middleware"
	userRepository "supermarket-backend/internal/repository/user"
	"supermarket-backend/internal/routes"
	authService "supermarket-backend/internal/service/auth"
	"sync/atomic"
	"syscall"
	"time"

	_ "supermarket-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	cfg := config.Load()

	// Setup signal context
	rootCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Database
	database, err := db.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL!")

	// Infrastructure
	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	// Repositories
	userRepo := userRepository.NewRepository(database)

	// Services
	authSvc := authService.NewService(
		userRepo,
		jwtService,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)

	// Gin
	router := gin.New()

	// Gin middlewares
	router.Use(
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins:     cfg.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: false,
		}),
	)

	// Readiness endpoint
	router.GET("/healthz", func(c *gin.Context) {
		if isShuttingDown.Load() {
			c.String(http.StatusServiceUnavailable, "Shutting down")
			return
		}

		c.String(http.StatusOK, "OK")
	})

	// Swagger
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	// Routes
	routes.RegisterAuthRoutes(
		router,
		authHandler,
		middleware.Auth(jwtService),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server starting on :8080.")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for shutdown signal
	<-rootCtx.Done()

	stop()

	isShuttingDown.Store(true)

	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)

	log.Println(
		"Readiness check propagated, now waiting for ongoing requests to finish.",
	)

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		_shutdownPeriod,
	)
	defer cancel()

	err = server.Shutdown(shutdownCtx)

	if err != nil {
		log.Println(
			"Failed to wait for ongoing requests to finish, waiting for forced cancellation.",
		)

		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
