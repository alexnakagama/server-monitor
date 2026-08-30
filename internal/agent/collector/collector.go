package collector

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/model"
)

func CollectMetrics() (model.Metric, error) {
	cpu, err := CPUUsage()
	if err != nil {
		return model.Metric{}, err
	}

	memory, err := MemoryUsage()
	if err != nil {
		return model.Metric{}, err
	}

	disk, err := DiskUsage()
	if err != nil {
		return model.Metric{}, err
	}

	network, err := NetworkUsage()
	if err != nil {
		return model.Metric{}, err
	}

	return model.Metric{
		CPUUsage:       cpu,
		MemoryUsage:    memory,
		DiskUsage:      disk,
		NetworkReceive: network.Received,
		NetworkSent:    network.Sent,
		Timestamp:      time.Now(),
	}, nil
}
