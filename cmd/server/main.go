package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/alexnakagama/server-monitor/internal/auth"
	"github.com/alexnakagama/server-monitor/internal/config"
	"github.com/alexnakagama/server-monitor/internal/scheduler"
	"github.com/alexnakagama/server-monitor/internal/server/db"
	"github.com/alexnakagama/server-monitor/internal/server/handler"
	"github.com/alexnakagama/server-monitor/internal/server/repository"
	"github.com/alexnakagama/server-monitor/internal/server/service"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.NewPostgres(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	key, err := paseto.V4SymmetricKeyFromHex(cfg.PasetoKey)
	if err != nil {
		log.Fatal(err)
	}

	pasetoManager := auth.NewPasetoManager(key)

	healthHandler := handler.NewHealthHandler(database)

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository, pasetoManager)
	userHandler := handler.NewUserHandler(userService)

	serverRepository := repository.NewServerRepository(database)
	serverService := service.NewServerService(serverRepository)
	serverHandler := handler.NewServerHandler(serverService)

	agentRepository := repository.NewAgentRepository(database)
	agentService := service.NewAgentService(agentRepository)
	agentHandler := handler.NewAgentHandler(agentService)

	metricRepository := repository.NewMetricRepository(database)
	metricService := service.NewMetricService(metricRepository, cfg.MetricRetentionDays)
	metricHandler := handler.NewMetricHandler(metricService)

	metricScheduler := scheduler.NewScheduler(metricService)

	go metricScheduler.Run(ctx)

	mux := http.NewServeMux()

	// public endpoints
	mux.HandleFunc("POST /users/register", userHandler.HandleRegister)
	mux.HandleFunc("POST /users/login", userHandler.HandleLogin)

	// endpoints for health of the app
	mux.HandleFunc("GET /health/live", healthHandler.HandleLive)
	mux.HandleFunc("GET /health/ready", healthHandler.HandleReady)

	// private endpoints
	// user endpoints
	mux.Handle("GET /users/profile/me", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(userHandler.HandleProfile),
	))

	mux.Handle("PUT /users/profile/me", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(userHandler.HandleUpdateProfile),
	))

	mux.Handle("DELETE /users/profile/me", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(userHandler.HandleDeleteProfile),
	))

	mux.Handle("PUT /users/profile/me/password", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(userHandler.HandleChangePassword),
	))

	// server endpoints
	mux.Handle("POST /servers", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(serverHandler.HandleCreate),
		),
	))

	mux.Handle("GET /servers/{name}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleGetByName),
	))

	mux.Handle("GET /servers", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleGetAll),
	))

	mux.Handle("GET /servers/os/{os}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleGetByOS),
	))

	mux.Handle("GET /servers/hostname/{hostname}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleGetByHostname),
	))

	mux.Handle("DELETE /servers/hostname/{hostname}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(serverHandler.HandleDeleteByHostname),
		),
	))

	mux.Handle("PUT /servers/hostname/{hostname}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(serverHandler.HandleUpdateByHostname),
		),
	))

	mux.Handle("PATCH /servers/hostname/{hostname}/name", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(serverHandler.HandleUpdateNameByHostname),
		),
	))

	mux.Handle("PATCH /servers/hostname/{hostname}/os", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(serverHandler.HandleUpdateOSByHostname),
		),
	))

	// agents endpoints
	mux.Handle("POST /servers/{serverID}/agents", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(agentHandler.HandleCreate),
		),
	))

	mux.Handle("GET /agents/{agentID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(agentHandler.HandleGetByID),
	))

	mux.Handle("DELETE /agents/{agentID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		auth.RequireRole(
			"admin",
			http.HandlerFunc(agentHandler.HandleDelete),
		),
	))

	mux.Handle("POST /agents/metrics", auth.AgentAuthMiddleware(
		agentService,
		http.HandlerFunc(metricHandler.HandleCreate),
	))

	// metric endpoints
	mux.Handle("GET /metrics/{metricID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(metricHandler.HandleGetByID),
	))

	mux.Handle("GET /metrics/server/{serverID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(metricHandler.HandleGetByServerID),
	))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server running on port: 8080")

	go func() {
		<-ctx.Done()

		log.Println("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("server shutdown error: %v", err)
			return
		}

		log.Println("server stopped")
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
