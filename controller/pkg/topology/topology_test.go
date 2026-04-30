package topology

import (
	"testing"
)

func TestNewTopologyService(t *testing.T) {
	ts := NewTopologyService()
	if ts == nil {
		t.Fatal("NewTopologyService returned nil")
	}
}

func TestRegisterDevice(t *testing.T) {
	ts := NewTopologyService()

	err := ts.RegisterDevice("leaf1", "10.0.1.1", "leaf")
	if err != nil {
		t.Fatalf("RegisterDevice failed: %v", err)
	}

	dev := ts.GetDevice("leaf1")
	if dev == nil {
		t.Fatal("Device not found after registration")
	}

	if dev.ID != "leaf1" || dev.Address != "10.0.1.1" || dev.Role != "leaf" {
		t.Fatal("Device attributes don't match")
	}
}

func TestRegisterDuplicateDevice(t *testing.T) {
	ts := NewTopologyService()

	ts.RegisterDevice("leaf1", "10.0.1.1", "leaf")
	err := ts.RegisterDevice("leaf1", "10.0.1.1", "leaf")

	if err == nil {
		t.Fatal("Expected error on duplicate registration")
	}
}

func TestListDevices(t *testing.T) {
	ts := NewTopologyService()

	ts.RegisterDevice("leaf1", "10.0.1.1", "leaf")
	ts.RegisterDevice("leaf2", "10.0.2.1", "leaf")
	ts.RegisterDevice("spine1", "10.1.1.1", "spine")

	devices := ts.ListDevices()
	if len(devices) != 3 {
		t.Fatalf("Expected 3 devices, got %d", len(devices))
	}
}

func TestTopologySummary(t *testing.T) {
	ts := NewTopologyService()

	ts.RegisterDevice("leaf1", "10.0.1.1", "leaf")
	ts.RegisterDevice("spine1", "10.1.1.1", "spine")

	summary := ts.Summary()
	if summary == "" {
		t.Fatal("Summary is empty")
	}

	expected := "Topology: 2 devices, 0 links"
	if summary != expected {
		t.Fatalf("Expected %q, got %q", expected, summary)
	}
}
