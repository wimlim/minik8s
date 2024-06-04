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
			m.AddRule(serviceSpec.ClusterIP, uint16(port.Port), podIP, uint16(port.TargetPort))
		}
	}
}

func (m *IPVSManager) DeleteService(serviceSpec apiobj.ServiceSpec) {
	for _, port := range serviceSpec.Ports {
		svc := &ipvs.Service{
			Address:       net.ParseIP(serviceSpec.ClusterIP),
			Port:          uint16(port.Port),
			Protocol:      syscall.IPPROTO_TCP,
			AddressFamily: syscall.AF_INET,
		}

		if err := m.handle.DelService(svc); err != nil {
			fmt.Printf("Failed to delete IPVS service on port %d: %v\n", port.Port, err)
		}
	}
}

func (m *IPVSManager) AddRule(svcIP string, svcPort uint16, podIP string, podPort uint16) {
	svc := &ipvs.Service{
		Address:       net.ParseIP(svcIP),
		Port:          svcPort,
		Protocol:      syscall.IPPROTO_TCP,
		AddressFamily: syscall.AF_INET,
	}

	dest := &ipvs.Destination{
		Address: net.ParseIP(podIP),
		Port:    podPort,
		Weight:  1,
	}

	if err := m.handle.NewDestination(svc, dest); err != nil {
		fmt.Printf("Failed to add pod IP %s to IPVS service on port %c: %v\n", podIP, svcPort, err)
	}
}

func (m *IPVSManager) DeleteRule(svcIP string, svcPort uint16, podIP string, podPort uint16) {
	svc := &ipvs.Service{
		Address:       net.ParseIP(svcIP),
		Port:          svcPort,
		Protocol:      syscall.IPPROTO_TCP,
		AddressFamily: syscall.AF_INET,
	}

	dest := &ipvs.Destination{
		Address: net.ParseIP(podIP),
		Port:    podPort,
	}

	if err := m.handle.DelDestination(svc, dest); err != nil {
		fmt.Printf("Failed to delete pod IP %s from IPVS service on port %c: %v\n", podIP, svcPort, err)
	}
}

func (m *IPVSManager) Close() {
	if m.handle != nil {
		m.handle.Close()
	}
}
