package loadaverage

type Result struct {
	OneMin     float64
	FiveMin    float64
	FifteenMin float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewLoadAverageCollector() (CollectorI, error) {
	if ErrNotImplemented != nil {
		return nil, ErrNotImplemented
	}

	return &Collector{}, nil
}

var ErrNotImplemented error