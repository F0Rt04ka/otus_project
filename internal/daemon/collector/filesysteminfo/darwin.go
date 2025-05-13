//go:build darwin

package filesysteminfo

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ❯ df -m -i -l
// Filesystem     1M-blocks   Used Available Capacity iused      ifree %iused  Mounted on
// /dev/disk3s3s1    948584  10719    633858     2%  424963 4292946590    0%   /
// devfs                  0      0         0   100%     688          0  100%   /dev
// /dev/disk3s6      948584   4096    633858     1%       4 6490708160    0%   /System/Volumes/VM
// /dev/disk3s4      948584   6794    633858     2%    1270 6490708160    0%   /System/Volumes/Preboot
// /dev/disk3s2      948584     48    633858     1%      59 6490708160    0%   /System/Volumes/Update
// /dev/disk1s2         500      6       480     2%       1    4922680    0%   /System/Volumes/xarts
// /dev/disk1s1         500      5       480     2%      34    4922680    0%   /System/Volumes/iSCPreboot
// /dev/disk1s3         500      3       480     1%      56    4922680    0%   /System/Volumes/Hardware
// /dev/disk3s1      948584 291873    633858    32% 3507977 6490708160    0%   /System/Volumes/Data
// map auto_home          0      0         0   100%       0          0     -   /System/Volumes/Data/home
type Collector struct{}

func (c *Collector) Collect(result Result) error {
	dfCmd := exec.Command("df", "-m", "-i", "-l")
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
		used, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing used: %w", err)
		}
		usedPcent, err := strconv.ParseFloat(strings.TrimSuffix(fields[4], "%"), 32)
		if err != nil {
			return fmt.Errorf("error parsing used percent: %w", err)
		}
		iused, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing inodes used: %w", err)
		}
		ipcent, err := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(fields[6], "%"), "-", "0", 1), 32)
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
