//go:build !linux

package diskload

import (
	"fmt"
	"runtime"
)

type Collector struct{}

func (*Collector) Collect(result *Result) error {
	return fmt.Errorf("CPU usage not implemented for %s", runtime.GOOS)
}
