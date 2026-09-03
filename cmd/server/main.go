package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"aidanwoods.dev/go-paseto"

	"github.com/alexnakagama/server-monitor/internal/auth"
	"github.com/alexnakagama/server-monitor/internal/server/db"
	"github.com/alexnakagama/server-monitor/internal/server/handler"
	"github.com/alexnakagama/server-monitor/internal/server/repository"
	"github.com/alexnakagama/server-monitor/internal/server/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
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

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository, pasetoManager)
	userHandler := handler.NewUserHandler(userService)

	serverRepository := repository.NewServerRepository(database)
	serverService := service.NewServerService(serverRepository)
	serverHandler := handler.NewServerHandler(serverService)

	mux := http.NewServeMux()

	// the only public endpoints, register and login
	mux.HandleFunc("POST /users/register", userHandler.HandleRegister)
	mux.HandleFunc("POST /users/login", userHandler.HandleLogin)

	// private endpoints
	mux.Handle("GET /users/profile/me", auth.AuthMiddleware(
		pasetoManager,
		http.HandlerFunc(userHandler.HandleProfile),
	))

	mux.Handle("POST /servers", auth.AuthMiddleware(
		pasetoManager,
		http.HandlerFunc(serverHandler.HandleCreate),
	))

	mux.Handle("GET /servers/{name}", auth.AuthMiddleware(
		pasetoManager,
		http.HandlerFunc(serverHandler.HandleGetByName),
	))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("server running on port: 8080")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
