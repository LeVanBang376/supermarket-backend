package main

// @title SmallSupermarket APIs
// @version 1.0
// @description SmallSupermarket APIs
// @host localhost:8080
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token

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
	branchRepository "supermarket-backend/internal/repository/branch"
	brandRepository "supermarket-backend/internal/repository/brand"
	positionRepository "supermarket-backend/internal/repository/position"
	roleRepository "supermarket-backend/internal/repository/role"
	skuRepository "supermarket-backend/internal/repository/sku"
	stockRepository "supermarket-backend/internal/repository/stock"
	unitRepository "supermarket-backend/internal/repository/unit"
	userRepository "supermarket-backend/internal/repository/user"
	"supermarket-backend/internal/routes"
	authService "supermarket-backend/internal/service/auth"
	branchService "supermarket-backend/internal/service/branch"
	brandService "supermarket-backend/internal/service/brand"
	positionService "supermarket-backend/internal/service/position"
	roleService "supermarket-backend/internal/service/role"
	skuService "supermarket-backend/internal/service/sku"
	stockService "supermarket-backend/internal/service/stock"
	unitService "supermarket-backend/internal/service/unit"
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
	authMiddleware := middleware.Auth(jwtService)

	// Repositories
	userRepo := userRepository.NewRepository()
	branchRepo := branchRepository.NewRepository()
	brandRepo := brandRepository.NewRepository()
	unitRepo := unitRepository.NewRepository()
	skuRepo := skuRepository.NewRepository()
	stockRepo := stockRepository.NewRepository()
	roleRepo := roleRepository.NewRepository()
	positionRepo := positionRepository.NewRepository()

	// Services
	authSvc := authService.NewService(
		database,
		userRepo,
		jwtService,
	)

	branchSvc := branchService.NewService(
		database,
		branchRepo,
	)

	brandSvc := brandService.NewService(
		database,
		brandRepo,
	)

	unitSvc := unitService.NewService(
		database,
		unitRepo,
	)

	skuSvc := skuService.NewService(
		database,
		skuRepo,
	)

	stockSvc := stockService.NewService(
		database,
		stockRepo,
	)

	roleSvc := roleService.NewService(
		database,
		roleRepo,
	)

	positionSvc := positionService.NewService(
		database,
		positionRepo,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	branchHandler := handler.NewBranchHandler(branchSvc)
	brandHandler := handler.NewBrandHandler(brandSvc)
	unitHandler := handler.NewUnitHandler(unitSvc)
	skuHandler := handler.NewSKUHandler(skuSvc)
	stockHandler := handler.NewStockHandler(stockSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	positionHandler := handler.NewPositionHandler(positionSvc)

	// Gin
	router := gin.New()

	// Gin middlewares
	router.Use(
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins:     cfg.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		}),
	)

	// Readiness endpoint
	router.GET("/healthz", func(c *gin.Context) {
		if isShuttingDown.Load() {
			c.String(
				http.StatusServiceUnavailable,
				"Shutting down",
			)
			return
		}

		c.String(
			http.StatusOK,
			"OK",
		)
	})

	// Swagger
	router.GET(
		"/swagger/*any",
		gin.WrapH(httpSwagger.WrapHandler),
	)

	// Routes
	routes.RegisterAuthRoutes(
		router,
		authHandler,
		authMiddleware,
	)

	routes.RegisterBranchRoutes(
		router,
		branchHandler,
		authMiddleware,
	)

	routes.RegisterBrandRoutes(
		router,
		brandHandler,
		authMiddleware,
	)

	routes.RegisterUnitRoutes(
		router,
		unitHandler,
		authMiddleware,
	)

	routes.RegisterSKURoutes(
		router,
		skuHandler,
		authMiddleware,
	)

	routes.RegisterStockRoutes(
		router,
		stockHandler,
		authMiddleware,
	)

	routes.RegisterRoleRoutes(
		router,
		roleHandler,
		authMiddleware,
	)

	routes.RegisterPositionRoutes(
		router,
		positionHandler,
		authMiddleware,
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

	log.Println(
		"Received shutdown signal, shutting down.",
	)

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
