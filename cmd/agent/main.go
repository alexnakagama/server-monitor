package main

import (
	"log"
	"os"
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	apiURL := os.Getenv("API_URL")
	token := os.Getenv("AGENT_TOKEN")

	client := agent.NewClient(apiURL, token)

	agent := agent.New(100*time.Second, client)
	agent.Run()
}
