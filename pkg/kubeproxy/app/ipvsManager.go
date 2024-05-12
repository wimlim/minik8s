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

func (m *IPVSManager) UpdateService(serviceSpec apiobj.ServiceSpec, newPodIPs []string) {
	for _, port := range serviceSpec.Ports {
		svc := &ipvs.Service{
			Address:       net.ParseIP(serviceSpec.ClusterIP),
			Port:          uint16(port.Port),
			Protocol:      syscall.IPPROTO_TCP,
			AddressFamily: syscall.AF_INET,
		}

		existingService, err := m.handle.GetService(svc)
		if err != nil {
			fmt.Printf("Failed to get IPVS service for update on port %d: %v\n", port.Port, err)
			continue
		}

		existingDests, err := m.handle.GetDestinations(existingService)
		if err != nil {
			fmt.Printf("Failed to get destinations for IPVS service on port %d: %v\n", port.Port, err)
			continue
		}

		existingDestMap := make(map[string]*ipvs.Destination)
		for _, dest := range existingDests {
			ip := dest.Address.String()
			existingDestMap[ip] = dest
		}

		for _, podIP := range newPodIPs {
			dest := &ipvs.Destination{
				Address: net.ParseIP(podIP),
				Port:    uint16(port.TargetPort),
				Weight:  1,
			}
			if _, exists := existingDestMap[podIP]; !exists {
				if err := m.handle.NewDestination(existingService, dest); err != nil {
					fmt.Printf("Failed to add new pod IP %s to IPVS service on port %d: %v\n", podIP, port.Port, err)
				}
			}
		}

		for ip, dest := range existingDestMap {
			if !contains(newPodIPs, ip) {
				if err := m.handle.DelDestination(existingService, dest); err != nil {
					fmt.Printf("Failed to delete pod IP %s from IPVS service on port %d: %v\n", ip, port.Port, err)
				}
			}
		}
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, exists := set[item]
	return exists
}

func (m *IPVSManager) Close() {
	if m.handle != nil {
		m.handle.Close()
	}
}
