package filesysteminfo

type Result map[string]*FileSystemUsage

type FileSystemUsage struct {
	Path            string
	UsedMB          float64
	UsedPcent       float64
	UsedInodes      float64
	UsedInodesPcent float64
}

type CollectorI interface {
	Collect(result Result) error
}

func NewFilesystemInfoCollector() CollectorI {
	return &Collector{}
}
