package topology

import (
	"fmt"
	"sync"
)

/* Device in the network graph */
type Device struct {
	ID        string
	Address   string
	Role      string
	Reachable bool
}

/* Link between devices */
type Link struct {
	SourceID string
	DestID   string
	Status   string
}

/* Topology graph */
type TopologyGraph struct {
	mu      sync.RWMutex
	devices map[string]*Device
	links   []Link
}

func NewTopologyGraph() *TopologyGraph {
	return &TopologyGraph{
		devices: make(map[string]*Device),
		links:   make([]Link, 0),
	}
}

/* Add device to topology */
func (tg *TopologyGraph) AddDevice(id, address, role string) error {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	if _, exists := tg.devices[id]; exists {
		return fmt.Errorf("device %s already exists", id)
	}

	tg.devices[id] = &Device{
		ID:        id,
		Address:   address,
		Role:      role,
		Reachable: true,
	}

	return nil
}

/* Get device by ID */
func (tg *TopologyGraph) GetDevice(id string) *Device {
	tg.mu.RLock()
	defer tg.mu.RUnlock()

	return tg.devices[id]
}

/* List all devices */
func (tg *TopologyGraph) ListDevices() []*Device {
	tg.mu.RLock()
	defer tg.mu.RUnlock()

	devices := make([]*Device, 0, len(tg.devices))
	for _, dev := range tg.devices {
		devices = append(devices, dev)
	}
	return devices
}

/* Add link between devices */
func (tg *TopologyGraph) AddLink(sourceID, destID string) error {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	if tg.devices[sourceID] == nil || tg.devices[destID] == nil {
		return fmt.Errorf("one or both devices not found")
	}

	tg.links = append(tg.links, Link{
		SourceID: sourceID,
		DestID:   destID,
		Status:   "up",
	})

	return nil
}

/* Get all links */
func (tg *TopologyGraph) GetLinks() []Link {
	tg.mu.RLock()
	defer tg.mu.RUnlock()

	return tg.links
}

/* Check if path exists between two devices (simple BFS) */
func (tg *TopologyGraph) HasPath(sourceID, destID string) bool {
	tg.mu.RLock()
	defer tg.mu.RUnlock()

	visited := make(map[string]bool)
	queue := []string{sourceID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == destID {
			return true
		}

		if visited[current] {
			continue
		}
		visited[current] = true

		/* Find neighbors */
		for _, link := range tg.links {
			if link.SourceID == current && !visited[link.DestID] {
				queue = append(queue, link.DestID)
			}
		}
	}

	return false
}
