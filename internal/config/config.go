package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	Workers   int
	QueueSize int
}

func LoadConfig() Config {
	workers := getEnvIntOrDefault("WORKERS", 4, 1)       // min 1
	queueSize := getEnvIntOrDefault("QUEUE_SIZE", 64, 1) // min 1
	return Config{Workers: workers, QueueSize: queueSize}
}

func getEnvIntOrDefault(key string, defaultVal, minVal int) int {
	envValStr := os.Getenv(key)
	if envValStr == "" {
		return defaultVal
	}

	envValInt, err := strconv.Atoi(envValStr)
	if err != nil {
		slog.Warn(
			"invalid env value, using default",
			slog.String("key", key),
			slog.String("value", envValStr),
			slog.Any("error", err),
		)
		return defaultVal
	}

	if envValInt < minVal {
		slog.Warn(
			"env value below min, using default",
			slog.String("key", key),
			slog.Int("value", envValInt),
			slog.Int("min", minVal),
		)
		return defaultVal
	}

	return envValInt
}
