package loadaverage

type Result struct {
	OneMin     float64
	FiveMin    float64
	FifteenMin float64
}

type CollectorI interface {
	Collect(result *Result) error
}

func NewLoadAverageCollector() CollectorI {
	return &Collector{}
}
