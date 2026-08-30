package collector

import "syscall"

type diskStats struct {
	total uint64
	free  uint64
}

func readDiskStats() (diskStats, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		return diskStats{}, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)

	return diskStats{
		total: total,
		free:  free,
	}, nil
}

func DiskUsage() (float64, error) {
	disk, err := readDiskStats()
	if err != nil {
		return 0, err
	}

	used := disk.total - disk.free
	usage := float64(used) / float64(disk.total) * 100

	return usage, nil
}
