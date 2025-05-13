//go:build linux

package loadaverage

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Collector struct{}

func (c *Collector) Collect(result *Result) error {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return fmt.Errorf("failed to read /proc/loadavg: %w", err)
	}

	parts := strings.Fields(string(data))
	if len(parts) < 3 {
		return fmt.Errorf("unexpected format in /proc/loadavg")
	}

	// Парсим значения Load Average
	oneMin, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return fmt.Errorf("failed to parse 1-minute load average: %w", err)
	}

	fiveMin, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return fmt.Errorf("failed to parse 5-minute load average: %w", err)
	}

	fifteenMin, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return fmt.Errorf("failed to parse 15-minute load average: %w", err)
	}

	result.OneMin = oneMin
	result.FiveMin = fiveMin
	result.FifteenMin = fifteenMin

	return nil
}
