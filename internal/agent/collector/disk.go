package collector

type diskStats struct {
	total uint64
	free  uint64
}

func readDiskStats() (float64, error)

func DiskUsage() (float64, error)
