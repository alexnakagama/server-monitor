package collector

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type networkStats struct {
	Received uint64
	Sent     uint64
}

func readNetworkStats() (networkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return networkStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var received uint64
	var sent uint64

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) < 10 {
			continue
		}

		iface := strings.TrimSuffix(fields[0], ":")

		if iface == "lo" {
			continue
		}

		rx, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return networkStats{}, err
		}

		tx, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return networkStats{}, err
		}

		received += rx
		sent += tx
	}

	if err := scanner.Err(); err != nil {
		return networkStats{}, err
	}

	return networkStats{
		Received: received,
		Sent:     sent,
	}, nil
}

func NetworkUsage() (networkStats, error) {
	first, err := readNetworkStats()
	if err != nil {
		return networkStats{}, err
	}

	time.Sleep(time.Second)

	second, err := readNetworkStats()
	if err != nil {
		return networkStats{}, err
	}

	received := second.Received - first.Received
	sent := second.Sent - first.Sent

	return networkStats{
		Received: received,
		Sent:     sent,
	}, nil
}
