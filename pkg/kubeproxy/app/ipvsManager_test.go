package kubeproxy

import (
	"net"
	"syscall"
	"testing"

	"github.com/moby/ipvs"
)

func TestBaseIPVS(t *testing.T) {
	handle, err := ipvs.New("")
	if err != nil {
		t.Errorf("Failed to initialize IPVS handle: %v", err)
		return
	}

	svc := &ipvs.Service{
		Address:       net.ParseIP("192.168.8.9"),
		Port:          uint16(80),
		Protocol:      syscall.IPPROTO_TCP,
		AddressFamily: syscall.AF_INET,
		SchedName:     "rr",
	}

	if err := handle.NewService(svc); err != nil {
		t.Errorf("Failed to add IPVS service: %v", err)
		return
	}

	dest := &ipvs.Destination{
		Address: net.ParseIP("10.32.1.1"),
		Port:    uint16(8080),
	}

	if err := handle.NewDestination(svc, dest); err != nil {
		t.Errorf("Failed to add IPVS destination: %v", err)
	}

	// clean up
	if err := handle.DelService(svc); err != nil {
		t.Errorf("Failed to delete IPVS service: %v", err)
	}
}
