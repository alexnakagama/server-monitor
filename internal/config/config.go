package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	PasetoKey           string
	MetricRetentionDays int
}

func Load() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pasetoKey := os.Getenv("PASETO_KEY")

	retentionDays, err := strconv.Atoi(os.Getenv("METRIC_RETENTION_DAYS"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		DatabaseURL:         databaseURL,
		PasetoKey:           pasetoKey,
		MetricRetentionDays: retentionDays,
	}, nil
}
