package daemon

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/F0Rt04ka/otus_project/internal/daemon/collector"
	"github.com/F0Rt04ka/otus_project/internal/daemon/config"
	daemon_server "github.com/F0Rt04ka/otus_project/internal/daemon/server"
	"google.golang.org/grpc"
)

type Daemon struct {
	runner  *collector.Runner
	results *collector.ResultMap
}

func (d *Daemon) GetResults() *collector.ResultMap {
	return d.results
}

func RunDaemon(cfg *config.CollectorsConfig) (*Daemon, error) {
	statsResults := collector.NewCollectorResultMap(
		cfg.SecondsSaveStats,
		time.Duration(cfg.ClearStatsSecondsInterval)*time.Second,
	)
	statsResults.RunClearDataHandler(time.Now().Unix())
	runner := collector.NewCollectorRunner(statsResults, cfg)
	err := runner.RunAll()
	if err != nil {
		return nil, fmt.Errorf("failed to run collector daemon: %w", err)
	}

	return &Daemon{
		runner:  runner,
		results: statsResults,
	}, nil
}

func RunGRPCServer(cfg *config.GrpcConfig, daemon *Daemon) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	daemon_server.Register(grpcServer, daemon.GetResults())

	go func() {
		log.Printf("gRPC server start on %s\n", lis.Addr().String())

		grpcServer.Serve(lis)
		log.Printf("gRPC server stopped")
	}()

	return grpcServer, nil
}
