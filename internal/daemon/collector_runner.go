package daemon

import (
	"fmt"
	"time"

	"github.com/F0Rt04ka/otus_project/config"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collectors"
)

type CollectorRunner struct {
	result            *CollectorResultMap
	cpuCollector      *collectors.CPUUsageCollector
	loadCollector     *collectors.LoadAverageCollector
	diskLoadCollector *collectors.DiskLoadCollector
	filesystemCollector *collectors.FilesystemInfoCollector
}

func NewCollectorRunner(
	result *CollectorResultMap,
	cpuCollector *collectors.CPUUsageCollector,
	loadCollector *collectors.LoadAverageCollector,
	diskLoadCollector *collectors.DiskLoadCollector,
	filesystemCollector *collectors.FilesystemInfoCollector,
) *CollectorRunner {
	return &CollectorRunner{
		result:            result,
		cpuCollector:      cpuCollector,
		loadCollector:     loadCollector,
		diskLoadCollector: diskLoadCollector,
		filesystemCollector: filesystemCollector,
	}
}

func Run() {
	collectorResult := NewCollectorResultMap()
	runner := NewCollectorRunner(
		collectorResult,
		collectors.NewCPUUsageCollector(),
		collectors.NewLoadAverageCollector(),
		collectors.NewDiskLoadCollector(),
		collectors.NewFilesystemInfoCollector(),
	)

	go func() {
		// горутина для очистки старых данных
		ticker := time.NewTicker(time.Duration(config.Cfg.ClearStatsSecondsInterval) * time.Second)
		defer ticker.Stop()
		oldestTime := time.Now().Unix() - int64(config.Cfg.SecondsSaveStats)

		for {
			t := <-ticker.C
			for i := oldestTime; i < t.Unix()-int64(config.Cfg.SecondsSaveStats); i++ {
				oldestTime = i + 1
				collectorResult.DeleteStatsForTime(i)
			}
		}
	}()

	runner.RunCpuCollector()
	runner.RunLoadCollector()
	runner.RunDiskLoadCollector()

	N := 5
	M := 15

	go func(calculatePeriod uint) {
		ticker := time.NewTicker(time.Duration(N) * time.Second)
		defer ticker.Stop()

		for {
			t := <-ticker.C
			printResults(collectorResult, t.Unix(), calculatePeriod)
		}
	}(uint(M))

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


func printResults(results *CollectorResultMap, unixTime int64, secondForAvg uint) {
	cpuStats := struct {
		UserMode   []float64
		SystemMode []float64
		Idle       []float64
	}{}
	loadStats := struct {
		OneMin     []float64
		FiveMin    []float64
		FifteenMin []float64
	}{}
	diskStats := struct {
		TPS       []float64
		ReadKBps  []float64
		WriteKBps []float64
	}{}

	for i := unixTime; i > unixTime-int64(secondForAvg); i-- {
		if res, _ := results.GetCPUStats(i); res != nil {
			cpuStats.UserMode = append(cpuStats.UserMode, res.UserMode)
			cpuStats.SystemMode = append(cpuStats.SystemMode, res.SystemMode)
			cpuStats.Idle = append(cpuStats.Idle, res.Idle)
		}
		if res, _ := results.GetLoadStats(i); res != nil {
			loadStats.OneMin = append(loadStats.OneMin, res.OneMin)
			loadStats.FiveMin = append(loadStats.FiveMin, res.FiveMin)
			loadStats.FifteenMin = append(loadStats.FifteenMin, res.FifteenMin)
		}
		if res, _ := results.GetDiskLoadStats(i); res != nil {
			diskStats.TPS = append(diskStats.TPS, res.TPS)
			diskStats.ReadKBps = append(diskStats.ReadKBps, res.ReadKBps)
			diskStats.WriteKBps = append(diskStats.WriteKBps, res.WriteKBps)
		}
	}

	fmt.Printf("CPU Usage: %.2f%% %.2f%% %.2f%% \n", avg(cpuStats.UserMode), avg(cpuStats.SystemMode), avg(cpuStats.Idle))
	fmt.Printf("Load Average: %.2f %.2f %.2f \n", avg(loadStats.OneMin), avg(loadStats.FiveMin), avg(loadStats.FifteenMin))
	fmt.Printf("Disk Load: %.2f TPS %.2f KB/s %.2f KB/s \n", avg(diskStats.TPS), avg(diskStats.ReadKBps), avg(diskStats.WriteKBps))
	// fmt.Println("Filesystem Usage:")
	// for _, fsInfo := range currentResult.FilesystemStats {
	// 	fmt.Printf("  %s: Used: %d MB (%.2f%%), Used Inodes: %d (%.2f%%)\n",
	// 		fsInfo.Path, fsInfo.UsedMB, fsInfo.UsedPcent, fsInfo.UsedInodes, fsInfo.UsedInodesPcent)
	// }
	fmt.Println("-----------------------------------------------------")
}

func avg(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	var sum float64
	for _, d := range data {
		sum += d
	}
	return sum / float64(len(data))
}
