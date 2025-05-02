package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	cfg Config
)

type Config struct {
	DebugMode        bool             // режим отладки
	GRPCConfig       GrpcConfig       `json:"grpc"`
	CollectorsConfig CollectorsConfig `json:"collectors"`
}

type CollectorsConfig struct {
	SecondsSaveStats          int  `json:"seconds_save_stats"`
	ClearStatsSecondsInterval int  `json:"clear_stats_seconds_interval"`
	EnableCPUUsage            bool `json:"enable_cpu_usage"`
	EnableLoadAverage         bool `json:"enable_load_average"`
	EnableDiskLoad            bool `json:"enable_disk_load"`
	EnableFilesystemInfo      bool `json:"enable_filesystem_info"`
}

type GrpcConfig struct {
	Port int `json:"port"`
}

func Load() (*Config, error) {
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.CollectorsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	cfg.DebugMode = true

	return &cfg, nil
}
