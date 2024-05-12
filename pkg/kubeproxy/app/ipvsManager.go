package kubeproxy

import (
	"fmt"
	"minik8s/pkg/apiobj"
	"net"
	"syscall"

	"github.com/moby/ipvs"
)

type IPVSManager struct {
	handle *ipvs.Handle
}

func NewIPVSManager() *IPVSManager {
	handle, err := ipvs.New("")
	if err != nil {
		fmt.Println("Failed to initialize IPVS handle:", err)
		return nil
	}
	return &IPVSManager{handle: handle}
}

func (m *IPVSManager) AddService(serviceSpec apiobj.ServiceSpec, podIPs []string) {
	for _, port := range serviceSpec.Ports {
		svc := &ipvs.Service{
			Address:       net.ParseIP(serviceSpec.ClusterIP),
			Port:          uint16(port.Port),
			Protocol:      syscall.IPPROTO_TCP,
			AddressFamily: syscall.AF_INET,
			SchedName:     "rr",
		}

		err := m.handle.NewService(svc)
		if err != nil {
			fmt.Printf("Failed to add IPVS service on port %d: %v\n", port.Port, err)
			continue
		}

		for _, podIP := range podIPs {
			dest := &ipvs.Destination{
				Address: net.ParseIP(podIP),
				Port:    uint16(port.TargetPort),
				Weight:  1,
			}
			if err := m.handle.NewDestination(svc, dest); err != nil {
				fmt.Printf("Failed to add pod IP %s to IPVS service on port %d: %v\n", podIP, port.Port, err)
			}
		}
	}
}

func (m *IPVSManager) DeleteService(serviceID string) {
}

func (m *IPVSManager) UpdateService(serviceSpec apiobj.ServiceSpec, podIPs []string) {
}

func (m *IPVSManager) Close() {
	if m.handle != nil {
		m.handle.Close()
	}
}
