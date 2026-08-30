package main

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/agent"
)

func main() {
	agent := agent.New(5 * time.Second)

	agent.Run()
}
