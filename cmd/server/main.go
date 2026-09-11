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
	"github.com/alexnakagama/server-monitor/internal/server/db"
	"github.com/alexnakagama/server-monitor/internal/server/handler"
	"github.com/alexnakagama/server-monitor/internal/server/repository"
	"github.com/alexnakagama/server-monitor/internal/server/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	database, err := db.NewPostgres(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	key, err := paseto.V4SymmetricKeyFromHex(os.Getenv("PASETO_KEY"))
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
	metricService := service.NewMetricService(metricRepository)
	metricHandler := handler.NewMetricHandler(metricService)

	mux := http.NewServeMux()

	// the only public endpoints, register and login
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
		http.HandlerFunc(serverHandler.HandleCreate),
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
		http.HandlerFunc(serverHandler.HandleDeleteByHostname),
	))

	mux.Handle("PUT /servers/hostname/{hostname}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleUpdateByHostname),
	))

	mux.Handle("PATCH /servers/hostname/{hostname}/name", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleUpdateNameByHostname),
	))

	mux.Handle("PATCH /servers/hostname/{hostname}/os", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(serverHandler.HandleUpdateOSByHostname),
	))

	// agents endpoints
	mux.Handle("POST /servers/{serverID}/agents", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(agentHandler.HandleCreate),
	))

	mux.Handle("GET /agents/{agentID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(agentHandler.HandleGetByID),
	))

	mux.Handle("DELETE /agents/{agentID}", auth.AuthMiddleware(
		pasetoManager,
		userRepository,
		http.HandlerFunc(agentHandler.HandleDelete),
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

	// creates the channel
	// receives signals of the operating system
	shutdownSignal := make(chan os.Signal, 1)

	// function which connects the signals of the os to the go program
	// when the os sends os.Interrupt or syscall.SIGTERM send that signal through the channel
	signal.Notify(
		shutdownSignal,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		// here the go routine stays blocked until it receives a shutdown signal
		<-shutdownSignal

		log.Println("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// shutdowns the server
		err := server.Shutdown(ctx)
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
