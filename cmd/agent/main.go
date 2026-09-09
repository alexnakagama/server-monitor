package main

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent"
)

func main() {
	client := agent.NewClient(
		"http://localhost:8080",
		"TU_AGENT_TOKEN",
	)

	agent := agent.New(10*time.Second, client)
	agent.Run()
}
