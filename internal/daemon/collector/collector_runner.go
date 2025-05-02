package collector

import (
	"fmt"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon/collectors"
	"github.com/F0Rt04ka/otus_project/internal/daemon/config"
)

type CollectorRunner struct {
	result              *CollectorResultMap
	cpuCollector        *collectors.CPUUsageCollector
	loadCollector       *collectors.LoadAverageCollector
	diskLoadCollector   *collectors.DiskLoadCollector
	filesystemCollector *collectors.FilesystemInfoCollector
}

func NewCollectorRunner(
	result *CollectorResultMap,
	cfg *config.CollectorsConfig,
) *CollectorRunner {

	runner := &CollectorRunner{result: result}

	if cfg.EnableCPUUsage {
		runner.cpuCollector = collectors.NewCPUUsageCollector()
	}
	if cfg.EnableLoadAverage {
		runner.loadCollector = collectors.NewLoadAverageCollector()
	}
	if cfg.EnableDiskLoad {
		runner.diskLoadCollector = collectors.NewDiskLoadCollector()
	}
	if cfg.EnableFilesystemInfo {
		runner.filesystemCollector = collectors.NewFilesystemInfoCollector()
	}

	return runner
}

func (c *CollectorRunner) RunAll() error {
	if c.cpuCollector != nil {
		if err := c.RunCpuCollector(); err != nil {
			return err
		}
	}
	if c.loadCollector != nil {
		if err := c.RunLoadCollector(); err != nil {
			return err
		}
	}
	if c.diskLoadCollector != nil {
		if err := c.RunDiskLoadCollector(); err != nil {
			return err
		}
	}

	return nil
}

func (c *CollectorRunner) RunCpuCollector() error {
	if c.cpuCollector == nil {
		return fmt.Errorf("cpu collector is not initialized")
	}

	go func() {
		ticker := time.NewTicker(700 * time.Millisecond)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &collectors.CPUUsageResult{}
			err := c.cpuCollector.Collect(result)
			if err != nil {
				panic(err)
			}

			c.result.AddCPUStats(collectTime.Unix(), result)
		}
	}()

	return nil
}

func (c *CollectorRunner) RunLoadCollector() error {
	if c.loadCollector == nil {
		return fmt.Errorf("load collector is not initialized")
	}

	go func() {
		ticker := time.NewTicker(700 * time.Millisecond)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &collectors.LoadAverageResult{}
			err := c.loadCollector.Collect(result)
			if err != nil {
				panic(err)
			}

			c.result.AddLoadStats(collectTime.Unix(), result)
		}
	}()

	return nil
}

func (c *CollectorRunner) RunDiskLoadCollector() error {
	if c.diskLoadCollector == nil {
		return fmt.Errorf("disk load collector is not initialized")
	}

	go func() {
		ticker := time.NewTicker(1500 * time.Millisecond)
		defer ticker.Stop()

		for {
			collectTime := <-ticker.C
			result := &collectors.DiskLoadResult{}
			err := c.diskLoadCollector.Collect(result)
			if err != nil {
				panic(err)
			}

			c.result.AddDiskLoadStats(collectTime.Unix(), result)
		}
	}()

	return nil
}

func (c *CollectorRunner) RunFilesystemCollector() error {
	// if c.filesystemCollector == nil {
	// 	return fmt.Errorf("filesystem collector is not initialized")
	// }

	// go func() {
	// 	ticker := time.NewTicker(2000 * time.Millisecond)
	// 	defer ticker.Stop()

	// 	for {
	// 		collectTime := <-ticker.C
	// 		result := &collectors.FilesystemInfoResult{}
	// 		err := c.filesystemCollector.Collect(result)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		c.result.AddFilesystemStats(collectTime.Unix(), result)
	// 	}
	// }()
	// TODO

	return nil
}
