package diskload

type Result struct {
	TPS       float64
	ReadKBps  float64
	WriteKBps float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewDiskLoadCollector() CollectorI {
	return &Collector{}
}
