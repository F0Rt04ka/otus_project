package cpuusage

type Result struct {
	UserMode   float64
	SystemMode float64
	Idle       float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewCPUUsageCollector() CollectorI {
	return &Collector{}
}
