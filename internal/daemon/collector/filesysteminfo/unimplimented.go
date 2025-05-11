//go:build !linux && !darwin

package filesysteminfo

import (
	"fmt"
	"runtime"
)

type Collector struct{}

func (*Collector) Collect(result Result) error {
	return ErrNotImplemented
}

func init() {
	ErrNotImplemented = fmt.Errorf("Filesystem info not implemented for %s", runtime.GOOS)
}