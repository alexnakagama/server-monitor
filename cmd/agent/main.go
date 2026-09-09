package main

import (
	"os"
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent"
)

func main() {
	apiURL := os.Getenv("API_URL")
	token := os.Getenv("AGENT_TOKEN")

	client := agent.NewClient(apiURL, token)

	agent := agent.New(100*time.Second, client)
	agent.Run()
}
