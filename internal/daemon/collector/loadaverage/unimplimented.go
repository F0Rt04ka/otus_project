//go:build !linux

package loadaverage

import (
	"fmt"
	"runtime"
)

type Collector struct{}

func (*Collector) Collect(result *Result) error {
	return fmt.Errorf("Filesystem info not implemented for %s", runtime.GOOS)
}
