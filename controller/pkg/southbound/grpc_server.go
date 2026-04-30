package southbound

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gundu/networking-sdn/controller/api"
	"github.com/gundu/networking-sdn/controller/pkg/topology"
	"google.golang.org/grpc"
)

type FabricAgentServer struct {
	api.UnimplementedFabricAgentServer
	topology *topology.TopologyService
}

func NewFabricAgentServer(ts *topology.TopologyService) *FabricAgentServer {
	return &FabricAgentServer{
		topology: ts,
	}
}

func (s *FabricAgentServer) RegisterDevice(ctx context.Context, info *api.DeviceInfo) (*api.DeviceID, error) {
	err := s.topology.RegisterDevice(info.DeviceId, info.DeviceAddr, info.DeviceRole)

	if err != nil {
		return &api.DeviceID{
			Id:         info.DeviceId,
			Registered: false,
		}, err
	}

	return &api.DeviceID{
		Id:         info.DeviceId,
		Registered: true,
	}, nil
}

func (s *FabricAgentServer) GetDeviceState(ctx context.Context, id *api.DeviceID) (*api.DeviceState, error) {
	device := s.topology.GetDevice(id.Id)
	if device == nil {
		return nil, fmt.Errorf("device not found: %s", id.Id)
	}

	return &api.DeviceState{
		DeviceId:          device.ID,
		Routes:            make([]*api.RouteInfo, 0),
		Tunnels:           make([]*api.TunnelInfo, 0),
		PacketsForwarded:  0,
		PacketsDropped:    0,
	}, nil
}

func (s *FabricAgentServer) CreateVxlanTunnel(ctx context.Context, config *api.TunnelConfig) (*api.TunnelStatus, error) {
	return &api.TunnelStatus{
		TunnelId: config.TunnelId,
		Created:  true,
		Status:   "active",
	}, nil
}

func (s *FabricAgentServer) AdvertiseBgpRoute(ctx context.Context, route *api.RouteAdvertisement) (*api.RouteStatus, error) {
	return &api.RouteStatus{
		Advertised:    true,
		PeersReceived: 0,
	}, nil
}

func (s *FabricAgentServer) ApplyAcl(ctx context.Context, rule *api.AclRule) (*api.AclStatus, error) {
	return &api.AclStatus{
		RuleId:  rule.RuleId,
		Applied: true,
	}, nil
}

func (s *FabricAgentServer) StreamDeviceEvents(id *api.DeviceID, stream api.FabricAgent_StreamDeviceEventsServer) error {
	fmt.Printf("Event stream started for device: %s\n", id.Id)
	return nil
}

type GrpcServer struct {
	server *grpc.Server
	lis    net.Listener
}

func NewGrpcServer(addr string, ts *topology.TopologyService) (*GrpcServer, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	agent := NewFabricAgentServer(ts)
	api.RegisterFabricAgentServer(grpcServer, agent)

	return &GrpcServer{
		server: grpcServer,
		lis:    lis,
	}, nil
}

func (gs *GrpcServer) Start() error {
	fmt.Printf("gRPC server listening on %s\n", gs.lis.Addr())
	return gs.server.Serve(gs.lis)
}

func (gs *GrpcServer) Stop() {
	gs.server.GracefulStop()
	gs.lis.Close()
}
