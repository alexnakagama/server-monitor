package collector

type diskStats struct {
	total uint64
	free  uint64
}

func DiskUsage() (float64, error)
