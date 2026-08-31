package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userHandler.HandleRegister)

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
