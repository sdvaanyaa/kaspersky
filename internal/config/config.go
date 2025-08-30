package config

import (
	"log/slog"
	"os"
	"strconv"
)

const (
	DefaultWorkers = 4
	MinWorkers     = 1
	MaxWorkers     = 32

	DefaultQueueSize = 64
	MinQueueSize     = 1
	MaxQueueSize     = 1024
)

type Config struct {
	Workers   int
	QueueSize int
}

func LoadConfig() Config {
	workers := getEnvIntOrDefault(
		"WORKERS",
		DefaultWorkers,
		MinWorkers,
		MaxWorkers,
	)
	queueSize := getEnvIntOrDefault(
		"QUEUE_SIZE",
		DefaultQueueSize,
		MinQueueSize,
		MaxQueueSize,
	)

	return Config{Workers: workers, QueueSize: queueSize}
}

func getEnvIntOrDefault(key string, defaultVal, minVal, maxVal int) int {
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
			"env value below minimum, using default",
			slog.String("key", key),
			slog.Int("value", envValInt),
			slog.Int("min", minVal),
		)
		return defaultVal
	}

	if envValInt > maxVal {
		slog.Warn(
			"env value above maximum, using default",
			slog.String("key", key),
			slog.Int("value", envValInt),
			slog.Int("max", maxVal),
		)
		return defaultVal
	}
	return envValInt
}
