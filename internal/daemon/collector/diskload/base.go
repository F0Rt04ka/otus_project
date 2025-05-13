package diskload

type Result struct {
	TPS       float64
	ReadKBps  float64
	WriteKBps float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewDiskLoadCollector() (CollectorI, error) {
	if ErrNotImplemented != nil {
		return nil, ErrNotImplemented
	}

	return &Collector{}, nil
}

var ErrNotImplemented error
