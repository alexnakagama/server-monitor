package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	PasetoKey           string
	MetricRetentionDays int
}

type AgentConfig struct {
	APIURL     string
	AgentToken string
	Interval   time.Duration
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

func LoadAgent() (AgentConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("no .env file found")
	}

	interval, err := time.ParseDuration(os.Getenv("AGENT_INTERVAL"))
	if err != nil {
		return AgentConfig{}, err
	}

	return AgentConfig{
		APIURL:     os.Getenv("API_URL"),
		AgentToken: os.Getenv("AGENT_TOKEN"),
		Interval:   interval,
	}, nil
}
