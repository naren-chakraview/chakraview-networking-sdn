package topology

import (
	"fmt"
	"sync"
)

/* Discovery service manages device registration */
type DiscoveryService struct {
	mu       sync.RWMutex
	devices  map[string]*Device
	handlers map[string][]DiscoveryHandler
}

type DiscoveryHandler func(event string, device *Device)

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{
		devices:  make(map[string]*Device),
		handlers: make(map[string][]DiscoveryHandler),
	}
}

/* Register a device */
func (ds *DiscoveryService) RegisterDevice(id, address, role string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if _, exists := ds.devices[id]; exists {
		return fmt.Errorf("device %s already registered", id)
	}

	device := &Device{
		ID:        id,
		Address:   address,
		Role:      role,
		Reachable: true,
	}

	ds.devices[id] = device

	/* Notify listeners */
	if handlers, ok := ds.handlers["device.registered"]; ok {
		for _, h := range handlers {
			go h("device.registered", device)
		}
	}

	return nil
}

/* Get registered device */
func (ds *DiscoveryService) GetDevice(id string) *Device {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.devices[id]
}

/* List all registered devices */
func (ds *DiscoveryService) ListDevices() []*Device {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	devices := make([]*Device, 0, len(ds.devices))
	for _, dev := range ds.devices {
		devices = append(devices, dev)
	}
	return devices
}

/* Mark device as unreachable */
func (ds *DiscoveryService) MarkUnreachable(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	device, exists := ds.devices[id]
	if !exists {
		return fmt.Errorf("device %s not found", id)
	}

	device.Reachable = false

	if handlers, ok := ds.handlers["device.unreachable"]; ok {
		for _, h := range handlers {
			go h("device.unreachable", device)
		}
	}

	return nil
}

/* Subscribe to discovery events */
func (ds *DiscoveryService) Subscribe(eventType string, handler DiscoveryHandler) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.handlers[eventType] = append(ds.handlers[eventType], handler)
}
