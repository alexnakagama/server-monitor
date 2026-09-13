package config

type Config struct {
	DatabaseURL         string
	PasetoKey           string
	MetricRetentionDays int
}

func Load() (Config, error) {}
