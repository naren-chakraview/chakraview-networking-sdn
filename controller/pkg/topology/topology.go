package topology

import (
	"fmt"
	"sync"
)

/* Topology service combines graph and discovery */
type TopologyService struct {
	graph     *TopologyGraph
	discovery *DiscoveryService
	mu        sync.RWMutex
}

func NewTopologyService() *TopologyService {
	return &TopologyService{
		graph:     NewTopologyGraph(),
		discovery: NewDiscoveryService(),
	}
}

/* Register a device (called by fabric nodes during startup) */
func (ts *TopologyService) RegisterDevice(id, address, role string) error {
	/* Register in discovery */
	err := ts.discovery.RegisterDevice(id, address, role)
	if err != nil {
		return err
	}

	/* Add to graph */
	return ts.graph.AddDevice(id, address, role)
}

/* Get device */
func (ts *TopologyService) GetDevice(id string) *Device {
	return ts.graph.GetDevice(id)
}

/* List all devices */
func (ts *TopologyService) ListDevices() []*Device {
	return ts.graph.ListDevices()
}

/* Verify connectivity */
func (ts *TopologyService) IsReachable(sourceID, destID string) bool {
	source := ts.graph.GetDevice(sourceID)
	dest := ts.graph.GetDevice(destID)

	if source == nil || dest == nil || !source.Reachable || !dest.Reachable {
		return false
	}

	return ts.graph.HasPath(sourceID, destID)
}

/* Get topology summary */
func (ts *TopologyService) Summary() string {
	devices := ts.ListDevices()
	links := ts.graph.GetLinks()

	return fmt.Sprintf("Topology: %d devices, %d links",
		len(devices), len(links))
}
