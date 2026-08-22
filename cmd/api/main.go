// @title Supermarket API
// @version 1.0
// @description Supermarket backend API
// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"supermarket-backend/internal/config"
	"supermarket-backend/internal/database"
	"supermarket-backend/internal/middleware"
	"sync/atomic"
	"syscall"
	"time"

	_ "supermarket-backend/docs"

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
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Readiness endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "OK")
	})

	// Database
	pool, err := database.NewPostgresPool(
		cfg.DatabaseURL,
	)

	if err != nil {
		panic(err)
	}

	defer pool.Close()
	fmt.Println("Connected to PostgreSQL!")

	mux := http.NewServeMux()

	mux.Handle(
		"/swagger/",
		httpSwagger.WrapHandler,
	)

	// APIs
	// authRepository := auth.NewRepository(pool)
	// authService := auth.NewService(authRepository)
	// authHandler := auth.NewHandler(authService)
	// auth.RegisterRoutes(mux, authHandler)

	// Apply middlewares - Recovery and CORS
	handler := middleware.Recovery(middleware.CORS(cfg.AllowedOrigins)(mux))

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Println("Server starting on :8080.")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
