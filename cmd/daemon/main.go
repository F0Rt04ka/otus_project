package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon"
	"github.com/F0Rt04ka/otus_project/internal/daemon/collector"
	"github.com/F0Rt04ka/otus_project/internal/daemon/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	flag.IntVar(&cfg.GRPCConfig.Port, "server-port", 44044, "gRPC server port")
	flag.Parse()

	daemonApp, err := daemon.RunDaemon(&cfg.CollectorsConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to run daemon: %v", err))
	}

	_, err = daemon.RunGRPCServer(&cfg.GRPCConfig, daemonApp)
	if err != nil {
		panic(fmt.Sprintf("failed to run gRPC server: %v", err))
	}

	if cfg.DebugMode {
		printResultsDebug(daemonApp.GetResults())
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func printResultsDebug(results *collector.ResultMap) {
	printStats := func(results *collector.ResultMap, unixTime int64, secondForAvg int64) {
		cpuStats := results.GetAvgCPUStats(unixTime, secondForAvg)
		loadStats := results.GetAvgLoadStats(unixTime, secondForAvg)
		diskStats := results.GetAvgDiskLoadStats(unixTime, secondForAvg)
		fsStats := results.GetAvgFilesystemStats(unixTime, secondForAvg)

		fmt.Println("")
		fmt.Printf("CPU Usage: %.2f%% %.2f%% %.2f%% \n", cpuStats.UserMode, cpuStats.SystemMode, cpuStats.Idle)
		fmt.Printf("Load Average: %.2f %.2f %.2f \n", loadStats.OneMin, loadStats.FiveMin, loadStats.FifteenMin)
		fmt.Printf("Disk Load: %.2f TPS %.2f KB/s %.2f KB/s \n", diskStats.TPS, diskStats.ReadKBps, diskStats.WriteKBps)
		fmt.Println("Filesystem Usage:")
		for _, fsStat := range fsStats {
			fmt.Printf(
				"  %s: %.2f used MB, %.2f%%; %.2f used inodes, %.2f%%\n",
				fsStat.Path,
				fsStat.UsedMB,
				fsStat.UsedPcent,
				fsStat.UsedInodes,
				fsStat.UsedInodesPcent,
			)
		}

		fmt.Println("-----------------------------------------------------")
	}

	log.Println("Printing stats every 5 seconds.")

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for t := range ticker.C {
			printStats(results, t.Unix(), 5)
		}
	}()
}
