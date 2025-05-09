// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: system_monitor.proto

package sysmon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SystemMonitor_GetStats_FullMethodName = "/system_monitor.SystemMonitor/GetStats"
)

// SystemMonitorClient is the client API for SystemMonitor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemMonitorClient interface {
	GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StatsResponse], error)
}

type systemMonitorClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemMonitorClient(cc grpc.ClientConnInterface) SystemMonitorClient {
	return &systemMonitorClient{cc}
}

func (c *systemMonitorClient) GetStats(ctx context.Context, in *StatsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StatsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &SystemMonitor_ServiceDesc.Streams[0], SystemMonitor_GetStats_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StatsRequest, StatsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SystemMonitor_GetStatsClient = grpc.ServerStreamingClient[StatsResponse]

// SystemMonitorServer is the server API for SystemMonitor service.
// All implementations must embed UnimplementedSystemMonitorServer
// for forward compatibility.
type SystemMonitorServer interface {
	GetStats(*StatsRequest, grpc.ServerStreamingServer[StatsResponse]) error
	mustEmbedUnimplementedSystemMonitorServer()
}

// UnimplementedSystemMonitorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSystemMonitorServer struct{}

func (UnimplementedSystemMonitorServer) GetStats(*StatsRequest, grpc.ServerStreamingServer[StatsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedSystemMonitorServer) mustEmbedUnimplementedSystemMonitorServer() {}
func (UnimplementedSystemMonitorServer) testEmbeddedByValue()                       {}

// UnsafeSystemMonitorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemMonitorServer will
// result in compilation errors.
type UnsafeSystemMonitorServer interface {
	mustEmbedUnimplementedSystemMonitorServer()
}

func RegisterSystemMonitorServer(s grpc.ServiceRegistrar, srv SystemMonitorServer) {
	// If the following call pancis, it indicates UnimplementedSystemMonitorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SystemMonitor_ServiceDesc, srv)
}

func _SystemMonitor_GetStats_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StatsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SystemMonitorServer).GetStats(m, &grpc.GenericServerStream[StatsRequest, StatsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type SystemMonitor_GetStatsServer = grpc.ServerStreamingServer[StatsResponse]

// SystemMonitor_ServiceDesc is the grpc.ServiceDesc for SystemMonitor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SystemMonitor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system_monitor.SystemMonitor",
	HandlerType: (*SystemMonitorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStats",
			Handler:       _SystemMonitor_GetStats_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "system_monitor.proto",
}
