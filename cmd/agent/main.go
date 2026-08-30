package main

import (
	"fmt"
	"log"

	"github.com/alexnakagama/server-monitor/internal/agent/collector"
)

func main() {
	cpu, err := collector.CPUUsage()
	if err != nil {
		log.Fatal(err)
	}

	memory, err := collector.MemoryUsage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CPU: %.2f%%\n", cpu)
	fmt.Println("Memory: %.2f%%\n", memory)
}
