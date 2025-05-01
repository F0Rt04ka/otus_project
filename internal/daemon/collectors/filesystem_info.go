package collectors

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type FilesystemInfoResult map[string]*FileSystemUsage

type FileSystemUsage struct {
	Path            string
	UsedMB          uint64
	UsedPcent       float32
	UsedInodes      uint64
	UsedInodesPcent float32
}

type FilesystemInfoCollector struct{}

func NewFilesystemInfoCollector() *FilesystemInfoCollector {
	return &FilesystemInfoCollector{}
}

func (c *FilesystemInfoCollector) Collect(result FilesystemInfoResult) error {
	dfInodesCmd := exec.Command("df", "--exclude-type=tmpfs", "--exclude-type=efivarfs", "-m", "--output=source,used,pcent,iused,ipcent")
	output, err := dfInodesCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run df: %w", err)
	}

	lines := strings.Split(string(output), "\n")

	for _, line := range lines[1:] { // Пропускаем заголовок
		fields := strings.Fields(line)

		if len(fields) < 5 {
			continue
		}

		source := fields[0]
		used, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing used: %w", err)
		}
		usedPcent, err := strconv.ParseFloat(strings.TrimSuffix(fields[2], "%"), 32)
		if err != nil {
			return fmt.Errorf("error parsing used percent: %w", err)
		}
		iused, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing inodes used: %w", err)
		}
		ipcent, err := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(fields[4], "%"), "-", "0", 1), 32)
		if err != nil {
			return fmt.Errorf("error parsing inodes percent: %w", err)
		}

		if result[source] == nil {
			result[source] = &FileSystemUsage{}
		}

		result[source].Path = source
		result[source].UsedMB = used
		result[source].UsedPcent = float32(usedPcent)
		result[source].UsedInodes = iused
		result[source].UsedInodesPcent = float32(ipcent)
	}

	return nil
}
