package main

import (
	"fmt"
	"log"

	"github.com/alexnakagama/server-monitor/internal/agent/collector"
)

func main() {
	metric, err := collector.CollectMetrics()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", metric)
}
