package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DebugMode        bool             // режим отладки
	GRPCConfig       GrpcConfig       `json:"grpc"`
	CollectorsConfig CollectorsConfig `json:"collectors"`
}

type CollectorsConfig struct {
	SecondsSaveStats          int  `json:"secondsSaveStats"`
	ClearStatsSecondsInterval int  `json:"clearStatsSecondsInterval"`
	EnableCPUUsage            bool `json:"enableCpuUsage"`
	CPUUsageIntervalMs        int  `json:"cpuUsageIntervalMs"`
	EnableLoadAverage         bool `json:"enableLoadAverage"`
	LoadAverageIntervalMs     int  `json:"loadAverageIntervalMs"`
	EnableDiskLoad            bool `json:"enableDiskLoad"`
	DiskLoadIntervalMs        int  `json:"diskLoadIntervalMs"`
	EnableFilesystemInfo      bool `json:"enableFilesystemInfo"`
	FilesystemInfoIntervalMs  int  `json:"filesystemInfoIntervalMs"`
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

	cfg := Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.CollectorsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	cfg.DebugMode = false

	return &cfg, nil
}
