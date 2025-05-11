package cpuusage

type Result struct {
	UserMode   float64
	SystemMode float64
	Idle       float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewCPUUsageCollector() (CollectorI, error) {
	if ErrNotImplemented != nil {
		return nil, ErrNotImplemented
	}

	return &Collector{}, nil
}

var ErrNotImplemented error
