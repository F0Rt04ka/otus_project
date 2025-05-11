//go:build linux

package filesysteminfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Collector struct{}

// ❯ df -m --exclude-type=tmpfs --exclude-type=efivarfs --output=source,used,pcent,iused,ipcent
// Filesystem       Used Use%   IUsed IUse%
// /dev/nvme0n1p2  26205  29%  453062    8%
// /dev/nvme0n1p1    185  21%     322    1%
// /dev/nvme0n1p4 162632  48% 2286533   10%
// /dev/nvme0n1p5      7   2%       0     -
// /dev/sda3      646610  73%  100359    1%
func (c *Collector) Collect(result Result) error {
	dfCmd := exec.Command(
		"df",
		"--exclude-type=tmpfs",
		"--exclude-type=efivarfs",
		"-m",
		"--output=source,used,pcent,iused,ipcent",
	)
	output, err := dfCmd.Output()
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
		result[source].UsedMB = float64(used)
		result[source].UsedPcent = usedPcent
		result[source].UsedInodes = float64(iused)
		result[source].UsedInodesPcent = ipcent
	}

	return nil
}
