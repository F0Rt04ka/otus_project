package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"

	sysmon "github.com/F0Rt04ka/otus_project/proto/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var port int
	N := 5
	M := 5
	flag.IntVar(&N, "N", N, "N - seconds timeout")
	flag.IntVar(&M, "M", M, "M - seconds for average")
	flag.IntVar(&port, "server-port", 44044, "gRPC server port")
	flag.Parse()

	grpcAddress := net.JoinHostPort("localhost", strconv.Itoa(port))

	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := sysmon.NewSystemMonitorClient(conn)

	if N < 1 || N > 500 {
		panic("N must be greater than 1 and less than 500")
	}

	if M < 1 || M > 500 {
		panic("M must be greater than 1 and less than 500")
	}

	//nolint:gosec
	stream, err := client.GetStats(context.Background(), &sysmon.StatsRequest{N: int32(N), M: int32(M)})
	if err != nil {
		panic(err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		printResults(resp)
	}
}

func printResults(results *sysmon.StatsResponse) {
	cpuStats := results.GetCpuUsage()
	loadAvg := results.GetLoadAverage()
	diskLoad := results.GetDiskLoad()
	fsStats := results.GetDiskStats()

	fmt.Println("")
	fmt.Printf("CPU Usage: %.2f%% %.2f%% %.2f%% \n", cpuStats.UserMode, cpuStats.SystemMode, cpuStats.Idle)
	fmt.Printf("Load Average: %.2f %.2f %.2f \n", loadAvg.OneMin, loadAvg.FiveMin, loadAvg.FifteenMin)
	fmt.Printf("Disk Load: %.2f TPS %.2f KB/s %.2f KB/s \n", diskLoad.Tps, diskLoad.ReadKbps, diskLoad.WriteKbps)
	fmt.Println("Filesystem Usage:")
	for _, fsStat := range fsStats {
		fmt.Printf(
			"  %s: %.2f used MB, %.2f%%; %.2f used inodes, %.2f%%\n",
			fsStat.Path,
			fsStat.UsedMb,
			fsStat.UsedPcent,
			fsStat.UsedInodes,
			fsStat.UsedInodesPcent,
		)
	}
	// fmt.Println("Filesystem Usage:")
	// for _, fsInfo := range currentResult.FilesystemStats {
	// 	fmt.Printf("  %s: Used: %d MB (%.2f%%), Used Inodes: %d (%.2f%%)\n",
	// 		fsInfo.Path, fsInfo.UsedMB, fsInfo.UsedPcent, fsInfo.UsedInodes, fsInfo.UsedInodesPcent)
	// }
	fmt.Println("-----------------------------------------------------")
}
