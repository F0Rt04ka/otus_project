package server

import (
	"fmt"
	"time"

	daemon "github.com/F0Rt04ka/otus_project/internal/daemon/collector"
	sysmon "github.com/F0Rt04ka/otus_project/proto/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SystemMonitorApi struct {
	sysmon.UnimplementedSystemMonitorServer
	StatsResults *daemon.CollectorResultMap
}

type SystemMonitor interface {
	GetStats(*sysmon.StatsRequest, grpc.ServerStreamingServer[sysmon.StatsResponse]) error
}

func Register(gRPCServer *grpc.Server, statsResultMap *daemon.CollectorResultMap) {
	sysmon.RegisterSystemMonitorServer(gRPCServer, &SystemMonitorApi{StatsResults: statsResultMap})
}

func (s *SystemMonitorApi) GetStats(req *sysmon.StatsRequest, stream grpc.ServerStreamingServer[sysmon.StatsResponse]) error {
	if req.N < 3 || req.N > 60 {
		return status.Error(codes.InvalidArgument, "N must be greater than 3 and less than 60")
	}
	if req.M < 3 || req.M > 120 {
		return status.Error(codes.InvalidArgument, "M must be greater than 3 and less than 120")
	}

	ticker := time.NewTicker(time.Duration(req.N) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			cpuStat := s.StatsResults.GetAvgCpuStats(t.Unix(), int64(req.M))

			resp := &sysmon.StatsResponse{
				CpuUsage: &sysmon.CPUUsageStat{
					UserMode:   cpuStat.UserMode,
					SystemMode: cpuStat.SystemMode,
					Idle:       cpuStat.Idle,
				},
			}

			if err := stream.Send(resp); err != nil {
				return fmt.Errorf("failed to send response: %w", err)
			}

		case <-stream.Context().Done():
			return nil
		}
	}
}
