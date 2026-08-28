package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type memoryStats struct {
	total     uint64
	available uint64
}

func readMemoryStats() (memoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return memoryStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var total uint64
	var available uint64

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			total, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memoryStats{}, err
			}

		case "MemAvailable:":
			available, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return memoryStats{}, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return memoryStats{}, err
	}

	return memoryStats{
		total:     total,
		available: available,
	}, nil
}

func MemoryUsage() (float64, error) {}
