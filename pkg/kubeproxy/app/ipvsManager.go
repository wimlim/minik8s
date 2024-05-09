package kubeproxy

import (
	"fmt"

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

func (m *IPVSManager) Close() {
	if m.handle != nil {
		m.handle.Close()
	}
}
