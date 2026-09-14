package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexnakagama/server-monitor/internal/agent"
	"github.com/alexnakagama/server-monitor/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.LoadAgent()
	if err != nil {
		log.Fatal(err)
	}

	client := agent.NewClient(cfg.APIURL, cfg.AgentToken)

	a := agent.New(cfg.Interval, client)

	log.Println("agent starting...")

	a.Run(ctx)

	log.Println("agent stopping...")
}
