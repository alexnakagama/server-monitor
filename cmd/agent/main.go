package main

import (
	"log"

	"github.com/alexnakagama/server-monitor/internal/agent"
	"github.com/alexnakagama/server-monitor/internal/config"
)

func main() {
	cfg, err := config.LoadAgent()
	if err != nil {
		log.Fatal(err)
	}

	client := agent.NewClient(cfg.APIURL, cfg.AgentToken)

	a := agent.New(cfg.Interval, client)
	a.Run()

	log.Println("agent stopping...")
}
