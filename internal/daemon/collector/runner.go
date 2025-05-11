package collector

import (
	"fmt"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/cpuusage"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/diskload"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/filesysteminfo"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector/loadaverage"
	"github.com/F0Rt04ka/otus_project/internal/daemon/config"
)

type Runner struct {
	errorChan                   chan error
	errorHandler                func(error)
	result                      *ResultMap
	cpuCollector                cpuusage.CollectorI
	cpuCollectorInterval        time.Duration
	loadCollector               loadaverage.CollectorI
	loadCollectorInterval       time.Duration
	diskLoadCollector           diskload.CollectorI
	diskLoadCollectorInterval   time.Duration
	filesystemCollector         filesysteminfo.CollectorI
	filesystemCollectorInterval time.Duration
}

func NewCollectorRunner(
	result *ResultMap,
	cfg *config.CollectorsConfig,
) *Runner {
	runner := &Runner{result: result}

	if cfg.EnableCPUUsage {
		runner.cpuCollector = cpuusage.NewCPUUsageCollector()
	}
	if cfg.EnableLoadAverage {
		runner.loadCollector = loadaverage.NewLoadAverageCollector()
	}
	if cfg.EnableDiskLoad {
		runner.diskLoadCollector = diskload.NewDiskLoadCollector()
	}
	if cfg.EnableFilesystemInfo {
		runner.filesystemCollector = filesysteminfo.NewFilesystemInfoCollector()
	}

	runner.cpuCollectorInterval = time.Duration(cfg.CPUUsageIntervalMs) * time.Millisecond
	runner.loadCollectorInterval = time.Duration(cfg.LoadAverageIntervalMs) * time.Millisecond
	runner.diskLoadCollectorInterval = time.Duration(cfg.DiskLoadIntervalMs) * time.Millisecond
	runner.filesystemCollectorInterval = time.Duration(cfg.FilesystemInfoIntervalMs) * time.Millisecond
	runner.errorChan = make(chan error)

	return runner
}

func (r *Runner) SetErrorHandler(handler func(error)) {
	r.errorHandler = handler
}

func (r *Runner) RunErrorHandler() {
	go func() {
		for err := range r.errorChan {
			if r.errorHandler != nil {
				r.errorHandler(err)
			} else {
				panic(err)
			}
		}
	}()
}

func (r *Runner) RunAll() {
	r.RunErrorHandler()

	if r.cpuCollector != nil {
		r.RunCPUCollector()
	}
	if r.loadCollector != nil {
		r.RunLoadCollector()
	}
	if r.diskLoadCollector != nil {
		r.RunDiskLoadCollector()
	}
	if r.filesystemCollector != nil {
		r.RunFilesystemCollector()
	}
}

func (r *Runner) RunCPUCollector() {
	if r.cpuCollector == nil {
		r.errorChan <- fmt.Errorf("cpu collector is not initialized")
		return
	}

	go func() {
		ticker := time.NewTicker(r.cpuCollectorInterval)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &cpuusage.Result{}
			err := r.cpuCollector.Collect(result)
			if err != nil {
				r.errorChan <- err
				continue
			}

			r.result.AddCPUStats(collectTime.Unix(), result)
		}
	}()
}

func (r *Runner) RunLoadCollector() {
	if r.loadCollector == nil {
		r.errorChan <- fmt.Errorf("load collector is not initialized")
		return
	}

	go func() {
		ticker := time.NewTicker(r.loadCollectorInterval)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &loadaverage.Result{}
			err := r.loadCollector.Collect(result)
			if err != nil {
				r.errorChan <- err
				continue
			}

			r.result.AddLoadStats(collectTime.Unix(), result)
		}
	}()
}

func (r *Runner) RunDiskLoadCollector() {
	if r.diskLoadCollector == nil {
		r.errorChan <- fmt.Errorf("disk load collector is not initialized")
		return
	}

	go func() {
		ticker := time.NewTicker(r.diskLoadCollectorInterval)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &diskload.Result{}
			err := r.diskLoadCollector.Collect(result)
			if err != nil {
				r.errorChan <- err
				continue
			}

			r.result.AddDiskLoadStats(collectTime.Unix(), result)
		}
	}()
}

func (r *Runner) RunFilesystemCollector() {
	if r.filesystemCollector == nil {
		r.errorChan <- fmt.Errorf("filesystem collector is not initialized")
		return
	}

	go func() {
		ticker := time.NewTicker(r.filesystemCollectorInterval)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := filesysteminfo.Result{}
			err := r.filesystemCollector.Collect(result)
			if err != nil {
				r.errorChan <- err
				continue
			}

			r.result.AddFilesystemStats(collectTime.Unix(), &result)
		}
	}()
}
