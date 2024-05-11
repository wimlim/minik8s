package kubeproxy

import (
	"fmt"
	"net"

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

func (m *IPVSManager) AddService(serviceDetails ipvs.Service) {
	err := m.handle.NewService(&serviceDetails)
	if err != nil {
		fmt.Println("Failed to add IPVS service:", err)
	}
}

func (m *IPVSManager) DeleteService(serviceDetails ipvs.Service) {
	err := m.handle.DelService(&serviceDetails)
	if err != nil {
		fmt.Println("Failed to delete IPVS service:", err)
	}
}

func (m *IPVSManager) UpdateService(serviceDetails ipvs.Service) {
	err := m.handle.UpdateService(&serviceDetails)
	if err != nil {
		fmt.Println("Failed to update IPVS service:", err)
	}
}

func (m *IPVSManager) AddPodToService(serviceIP string, podIP string) {
	dest := &ipvs.Destination{
		Address: net.ParseIP(podIP),
		Port:    0,
		Weight:  1,
	}
	svc := &ipvs.Service{
		Address: net.ParseIP(serviceIP),
		Port:    80,
	}

	if err := m.handle.NewDestination(svc, dest); err != nil {
		fmt.Println("Failed to add pod to IPVS service:", err)
	}
}

func (m *IPVSManager) AddPodsToService(serviceIP string, podIPs []string) {
	for _, podIP := range podIPs {
		m.AddPodToService(serviceIP, podIP)
	}
}

func (m *IPVSManager) Close() {
	if m.handle != nil {
		m.handle.Close()
	}
}
