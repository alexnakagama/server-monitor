package collector

type memoryStats struct {
	total     uint64
	available uint64
}

func readMemoryStats() (memoryStats, error) {}

func MemoryUsage() (float64, error) {}
