package collectors

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CPUUsageResult struct {
	UserMode   float64
	SystemMode float64
	Idle       float64
}

type CPUUsageCollector struct {
	prevUser   uint64
	prevSystem uint64
	prevIdle   uint64
	prevTotal  uint64
}

func NewCPUUsageCollector() *CPUUsageCollector {
	return &CPUUsageCollector{}
}

func (c *CPUUsageCollector) Collect(result *CPUUsageResult) error {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return fmt.Errorf("failed to read /proc/stat: %w", err)
	}

	// Ищем строку с общей статистикой CPU (начинается с "cpu ")
	lines := strings.Split(string(data), "\n")
	var cpuLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			cpuLine = line
			break
		}
	}

	if cpuLine == "" {
		return fmt.Errorf("failed to find CPU stats in /proc/stat")
	}

	fields := strings.Fields(cpuLine)
	if len(fields) < 5 {
		return fmt.Errorf("unexpected format in /proc/stat")
	}

	// Парсим значения времени
	// 1	user	Time spent with normal processing in user mode.
	// 2	nice	Time spent with niced processes in user mode.
	// 3	system	Time spent running in kernel mode.
	// 4	idle	Time spent in vacations twiddling thumbs.
	// 5	iowait	Time spent waiting for I/O to completed. This is considered idle time too. since 2.5.41
	// 6	irq	    Time spent serving hardware interrupts. since 2.6.0
	// 7	softirq	Time spent serving software interrupts. since 2.6.0
	// 8	steal	Time stolen by other operating systems running in a virtual environment. since 2.6.11
	// 9	guest	Time spent for running a virtual CPU or guest OS under the control of the kernel. since 2.6.24

	user, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse user time: %w", err)
	}

	system, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse system time: %w", err)
	}

	idle, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse idle time: %w", err)
	}

	total := user + system + idle

	// Если это первый вызов, сохраняем текущие значения и возвращаем пустой результат
	if c.prevTotal == 0 {
		c.prevUser = user
		c.prevSystem = system
		c.prevIdle = idle
		c.prevTotal = total
		return nil
	}

	deltaUser := user - c.prevUser
	deltaSystem := system - c.prevSystem
	deltaIdle := idle - c.prevIdle
	deltaTotal := total - c.prevTotal

	c.prevUser = user
	c.prevSystem = system
	c.prevIdle = idle
	c.prevTotal = total

	result.UserMode = float64(deltaUser) / float64(deltaTotal) * 100
	result.SystemMode = float64(deltaSystem) / float64(deltaTotal) * 100
	result.Idle = float64(deltaIdle) / float64(deltaTotal) * 100

	return nil
}
