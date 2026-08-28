package collector

type memoryStats struct {
	total     uint64
	available uint64
}

func MemoryUsage() (float64, error) {}
