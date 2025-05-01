package daemon

import (
	"sync"

	"github.com/F0Rt04ka/otus_project/config"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collectors"
)

type CollectorResultMap struct {
	cpuStats           map[int64]*collectors.CPUUsageResult
	cpuStatsMux        sync.RWMutex
	loadStats          map[int64]*collectors.LoadAverageResult
	loadStatsMux       sync.RWMutex
	diskLoadStats      map[int64]*collectors.DiskLoadResult
	diskLoadStatsMux   sync.RWMutex
	filesystemStats    map[int64]*collectors.FilesystemInfoResult
	filesystemStatsMux sync.Mutex
}

func NewCollectorResultMap() *CollectorResultMap {
	mapSize := config.Cfg.SecondsSaveStats + config.Cfg.ClearStatsSecondsInterval*2
	return &CollectorResultMap{
		cpuStats:        make(map[int64]*collectors.CPUUsageResult, mapSize),
		loadStats:       make(map[int64]*collectors.LoadAverageResult, mapSize),
		diskLoadStats:   make(map[int64]*collectors.DiskLoadResult, mapSize),
		filesystemStats: make(map[int64]*collectors.FilesystemInfoResult, mapSize),
	}
}

func (crm *CollectorResultMap) AddCPUStats(unixTime int64, result *collectors.CPUUsageResult) {
	crm.cpuStatsMux.Lock()
	defer crm.cpuStatsMux.Unlock()
	crm.cpuStats[unixTime] = result
}
func (crm *CollectorResultMap) GetCPUStats(unixTime int64) (*collectors.CPUUsageResult, bool) {
	crm.cpuStatsMux.RLock()
	defer crm.cpuStatsMux.RUnlock()
	result, exists := crm.cpuStats[unixTime]
	return result, exists
}

func (crm *CollectorResultMap) AddLoadStats(unixTime int64, result *collectors.LoadAverageResult) {
	crm.loadStatsMux.Lock()
	defer crm.loadStatsMux.Unlock()
	crm.loadStats[unixTime] = result
}
func (crm *CollectorResultMap) GetLoadStats(unixTime int64) (*collectors.LoadAverageResult, bool) {
	crm.loadStatsMux.Lock()
	defer crm.loadStatsMux.Unlock()
	result, exists := crm.loadStats[unixTime]
	return result, exists
}

func (crm *CollectorResultMap) AddDiskLoadStats(unixTime int64, result *collectors.DiskLoadResult) {
	crm.diskLoadStatsMux.Lock()
	defer crm.diskLoadStatsMux.Unlock()
	crm.diskLoadStats[unixTime] = result
}
func (crm *CollectorResultMap) GetDiskLoadStats(unixTime int64) (*collectors.DiskLoadResult, bool) {
	crm.diskLoadStatsMux.Lock()
	defer crm.diskLoadStatsMux.Unlock()
	result, exists := crm.diskLoadStats[unixTime]
	return result, exists
}

func (crm *CollectorResultMap) AddFilesystemStats(unixTime int64, result *collectors.FilesystemInfoResult) {
	crm.filesystemStatsMux.Lock()
	defer crm.filesystemStatsMux.Unlock()
	crm.filesystemStats[unixTime] = result
}
func (crm *CollectorResultMap) GetFilesystemStats(unixTime int64) (*collectors.FilesystemInfoResult, bool) {
	crm.filesystemStatsMux.Lock()
	defer crm.filesystemStatsMux.Unlock()
	result, exists := crm.filesystemStats[unixTime]
	return result, exists
}

func (crm *CollectorResultMap) DeleteStatsForTime(unixTime int64) {
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
