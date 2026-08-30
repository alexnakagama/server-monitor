package collector

type networkStats struct {
	received uint64
	sent     uint64
}

func readNetworkStats() (networkStats, error)
