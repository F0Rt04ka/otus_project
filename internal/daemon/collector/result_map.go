package collector

import (
	"math"
	"sync"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/cpuusage"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/diskload"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/filesysteminfo"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/loadaverage"
)

type ResultMap struct {
	secondsForSaveStats  int
	clearOldDataInterval time.Duration

	cpuStats           map[int64]*cpuusage.Result
	cpuStatsMux        sync.RWMutex
	loadStats          map[int64]*loadaverage.Result
	loadStatsMux       sync.RWMutex
	diskLoadStats      map[int64]*diskload.Result
	diskLoadStatsMux   sync.RWMutex
	filesystemStats    map[int64]*filesysteminfo.Result
	filesystemStatsMux sync.Mutex
}

func NewCollectorResultMap(secondsForSaveStats int, clearOldDataInterval time.Duration) *ResultMap {
	mapSize := secondsForSaveStats + int(clearOldDataInterval.Seconds()*2)

	return &ResultMap{
		secondsForSaveStats:  secondsForSaveStats,
		clearOldDataInterval: clearOldDataInterval,
		cpuStats:             make(map[int64]*cpuusage.Result, mapSize),
		loadStats:            make(map[int64]*loadaverage.Result, mapSize),
		diskLoadStats:        make(map[int64]*diskload.Result, mapSize),
		filesystemStats:      make(map[int64]*filesysteminfo.Result, mapSize),
	}
}

func (crm *ResultMap) RunClearDataHandler(startUnixTime int64) {
	go func() {
		// горутина для очистки старых данных
		ticker := time.NewTicker(crm.clearOldDataInterval)
		defer ticker.Stop()
		oldestTime := startUnixTime - int64(crm.secondsForSaveStats)

		for {
			t := <-ticker.C
			for i := oldestTime; i < t.Unix()-int64(crm.secondsForSaveStats); i++ {
				crm.DeleteStatsForTime(i)
				oldestTime = i + 1
			}
		}
	}()
}

func (crm *ResultMap) AddCPUStats(unixTime int64, result *cpuusage.Result) {
	crm.cpuStatsMux.Lock()
	defer crm.cpuStatsMux.Unlock()
	crm.cpuStats[unixTime] = result
}

func (crm *ResultMap) GetCPUStats(unixTime int64) (*cpuusage.Result, bool) {
	crm.cpuStatsMux.RLock()
	defer crm.cpuStatsMux.RUnlock()
	result, exists := crm.cpuStats[unixTime]
	return result, exists
}

func (crm *ResultMap) GetAvgCPUStats(unixTime int64, secondForAvg int64) *cpuusage.Result {
	stats := struct {
		UserMode   []float64
		SystemMode []float64
		Idle       []float64
	}{}

	for i := unixTime; i > unixTime-secondForAvg; i-- {
		if res, _ := crm.GetCPUStats(i); res != nil {
			stats.UserMode = append(stats.UserMode, res.UserMode)
			stats.SystemMode = append(stats.SystemMode, res.SystemMode)
			stats.Idle = append(stats.Idle, res.Idle)
		}
	}
	if len(stats.UserMode) == 0 {
		return nil
	}

	return &cpuusage.Result{
		UserMode:   avg(stats.UserMode),
		SystemMode: avg(stats.SystemMode),
		Idle:       avg(stats.Idle),
	}
}

func (crm *ResultMap) AddLoadStats(unixTime int64, result *loadaverage.Result) {
	crm.loadStatsMux.Lock()
	defer crm.loadStatsMux.Unlock()
	crm.loadStats[unixTime] = result
}

func (crm *ResultMap) GetLoadStats(unixTime int64) (*loadaverage.Result, bool) {
	crm.loadStatsMux.Lock()
	defer crm.loadStatsMux.Unlock()
	result, exists := crm.loadStats[unixTime]
	return result, exists
}

func (crm *ResultMap) GetAvgLoadStats(unixTime int64, secondForAvg int64) *loadaverage.Result {
	stats := struct {
		OneMin     []float64
		FiveMin    []float64
		FifteenMin []float64
	}{}

	for i := unixTime; i > unixTime-secondForAvg; i-- {
		if res, _ := crm.GetLoadStats(i); res != nil {
			stats.OneMin = append(stats.OneMin, res.OneMin)
			stats.FiveMin = append(stats.FiveMin, res.FiveMin)
			stats.FifteenMin = append(stats.FifteenMin, res.FifteenMin)
		}
	}

	return &loadaverage.Result{
		OneMin:     avg(stats.OneMin),
		FiveMin:    avg(stats.FiveMin),
		FifteenMin: avg(stats.FifteenMin),
	}
}

func (crm *ResultMap) AddDiskLoadStats(unixTime int64, result *diskload.Result) {
	crm.diskLoadStatsMux.Lock()
	defer crm.diskLoadStatsMux.Unlock()
	crm.diskLoadStats[unixTime] = result
}

func (crm *ResultMap) GetDiskLoadStats(unixTime int64) (*diskload.Result, bool) {
	crm.diskLoadStatsMux.Lock()
	defer crm.diskLoadStatsMux.Unlock()
	result, exists := crm.diskLoadStats[unixTime]
	return result, exists
}

func (crm *ResultMap) GetAvgDiskLoadStats(unixTime int64, secondForAvg int64) *diskload.Result {
	stats := struct {
		TPS       []float64
		ReadKBps  []float64
		WriteKBps []float64
	}{}

	for i := unixTime; i > unixTime-secondForAvg; i-- {
		if res, _ := crm.GetDiskLoadStats(i); res != nil {
			stats.TPS = append(stats.TPS, res.TPS)
			stats.ReadKBps = append(stats.ReadKBps, res.ReadKBps)
			stats.WriteKBps = append(stats.WriteKBps, res.WriteKBps)
		}
	}

	return &diskload.Result{
		TPS:       avg(stats.TPS),
		ReadKBps:  avg(stats.ReadKBps),
		WriteKBps: avg(stats.WriteKBps),
	}
}

func (crm *ResultMap) AddFilesystemStats(unixTime int64, result *filesysteminfo.Result) {
	crm.filesystemStatsMux.Lock()
	defer crm.filesystemStatsMux.Unlock()
	crm.filesystemStats[unixTime] = result
}

func (crm *ResultMap) GetFilesystemStats(unixTime int64) (*filesysteminfo.Result, bool) {
	crm.filesystemStatsMux.Lock()
	defer crm.filesystemStatsMux.Unlock()
	result, exists := crm.filesystemStats[unixTime]
	return result, exists
}

func (crm *ResultMap) GetAvgFilesystemStats(unixTime int64, secondForAvg int64) filesysteminfo.Result {
	stats := make(map[string]*struct {
		UsedMB          []float64
		UsedPcent       []float64
		UsedInodes      []float64
		UsedInodesPcent []float64
	})
	for i := unixTime; i > unixTime-secondForAvg; i-- {
		if res, _ := crm.GetFilesystemStats(i); res != nil {
			for _, fsInfo := range *res {
				if _, exists := stats[fsInfo.Path]; !exists {
					stats[fsInfo.Path] = &struct {
						UsedMB          []float64
						UsedPcent       []float64
						UsedInodes      []float64
						UsedInodesPcent []float64
					}{}
				}

				stats[fsInfo.Path].UsedMB = append(stats[fsInfo.Path].UsedMB, fsInfo.UsedMB)
				stats[fsInfo.Path].UsedPcent = append(stats[fsInfo.Path].UsedPcent, fsInfo.UsedPcent)
				stats[fsInfo.Path].UsedInodes = append(stats[fsInfo.Path].UsedInodes, fsInfo.UsedInodes)
				stats[fsInfo.Path].UsedInodesPcent = append(stats[fsInfo.Path].UsedInodesPcent, fsInfo.UsedInodesPcent)
			}
		}
	}
	if len(stats) == 0 {
		return nil
	}

	avgStats := make(filesysteminfo.Result, len(stats))

	for path, fsStat := range stats {
		avgStats[path] = &filesysteminfo.FileSystemUsage{
			Path:            path,
			UsedMB:          avg(fsStat.UsedMB),
			UsedPcent:       avg(fsStat.UsedPcent),
			UsedInodes:      avg(fsStat.UsedInodes),
			UsedInodesPcent: avg(fsStat.UsedInodesPcent),
		}
	}

	return avgStats
}

func (crm *ResultMap) DeleteStatsForTime(unixTime int64) {
	crm.cpuStatsMux.Lock()
	delete(crm.cpuStats, unixTime)
	crm.cpuStatsMux.Unlock()

	crm.loadStatsMux.Lock()
	delete(crm.loadStats, unixTime)
	crm.loadStatsMux.Unlock()

	crm.diskLoadStatsMux.Lock()
	delete(crm.diskLoadStats, unixTime)
	crm.diskLoadStatsMux.Unlock()

	crm.filesystemStatsMux.Lock()
	delete(crm.filesystemStats, unixTime)
	crm.filesystemStatsMux.Unlock()
}

func avg(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, d := range data {
		sum += d
	}

	return math.Round(sum/float64(len(data))*100) / 100
}
