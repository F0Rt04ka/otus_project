//go:build !linux

package diskload

import (
	"fmt"
	"runtime"
)

type Collector struct{}

func (*Collector) Collect(result *Result) error {
	return ErrNotImplemented
}

func init() {
	ErrNotImplemented = fmt.Errorf("Disk load not implemented for %s", runtime.GOOS)
}
