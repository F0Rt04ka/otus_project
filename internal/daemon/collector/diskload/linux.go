//go:build linux

package diskload

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Collector struct {
	prevReadSectors  uint64
	prevWriteSectors uint64
	prevIOs          uint64
	prevTime         int64
}

const sectorSize = 512 // Размер сектора в байтах

func (c *Collector) Collect(result *Result) error {
	data, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		return fmt.Errorf("failed to read /proc/diskstats: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var totalReadSectors, totalWriteSectors, totalIOs uint64

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}

		// Пропускаем устройства, которые не являются физическими дисками
		deviceName := fields[2]
		if !strings.HasPrefix(deviceName, "sd") && !strings.HasPrefix(deviceName, "nvme") {
			continue
		}

		// Читаем значения
		readSectors, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse read sectors: %w", err)
		}

		writeSectors, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse write sectors: %w", err)
		}

		ios, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse IOs: %w", err)
		}

		totalReadSectors += readSectors
		totalWriteSectors += writeSectors
		totalIOs += ios
	}

	currentTime := time.Now().Unix()

	// Если это первый вызов, сохраняем текущие значения и возвращаем пустой результат
	if c.prevTime == 0 {
		c.prevReadSectors = totalReadSectors
		c.prevWriteSectors = totalWriteSectors
		c.prevIOs = totalIOs
		c.prevTime = currentTime
		return nil
	}

	deltaReadSectors := totalReadSectors - c.prevReadSectors
	deltaWriteSectors := totalWriteSectors - c.prevWriteSectors
	deltaIOs := totalIOs - c.prevIOs
	deltaTime := float64(currentTime - c.prevTime)

	c.prevReadSectors = totalReadSectors
	c.prevWriteSectors = totalWriteSectors
	c.prevIOs = totalIOs
	c.prevTime = currentTime

	result.TPS = float64(deltaIOs) / deltaTime
	result.ReadKBps = float64(deltaReadSectors*sectorSize) / 1024.0 / deltaTime
	result.WriteKBps = float64(deltaWriteSectors*sectorSize) / 1024.0 / deltaTime

	return nil
}
