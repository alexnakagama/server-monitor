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

	disk, err := collector.DiskUsage()
	if err != nil {
		log.Fatal(err)
	}

	network, err := collector.NetworkUsage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CPU: %.2f%%\n", cpu)
	fmt.Printf("Memory: %.2f%%\n", memory)
	fmt.Printf("Disk: %.2f%%\n", disk)
	fmt.Printf("Nw recieved: %d bytes\n", network.Received)
	fmt.Printf("Nw sent: %d\n", network.Sent)
}
