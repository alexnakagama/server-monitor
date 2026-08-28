package collector

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpuStats struct {
	user    uint64
	nice    uint64
	system  uint64
	idle    uint64
	iowait  uint64
	irq     uint64
	softirq uint64
	steal   uint64
}

func readCPUStats() (cpuStats, error) {
	file, err := os.Open("proc/stat")
	if err != nil {
		return cpuStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return cpuStats{}, nil
	}

	fields := strings.Fields(scanner.Text())

	if len(fields) < 9 {
		return cpuStats{}, errors.New("invalid proc/stat format")
	}

	user, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	nice, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	system, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	idle, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	iowait, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	irq, err := strconv.ParseUint(fields[6], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	softirq, err := strconv.ParseUint(fields[7], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	steal, err := strconv.ParseUint(fields[8], 10, 64)
	if err != nil {
		return cpuStats{}, err
	}

	return cpuStats{
		user:    user,
		nice:    nice,
		system:  system,
		idle:    idle,
		iowait:  iowait,
		irq:     irq,
		softirq: softirq,
		steal:   steal,
	}, nil
}

func totalCPU(stats cpuStats) uint64 {}

func CPUUsage() (float64, error) {
	first, err := readCPUStats()
	if err != nil {
		return 0, err
	}

	time.Sleep(500 * time.Millisecond)

	second, err := readCPUStats()
	if err != nil {
		return 0, err
	}
}
