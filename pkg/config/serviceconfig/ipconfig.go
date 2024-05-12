package serviceconfig

import (
	"fmt"
)

const (
	BaseIp = "192.168.1."
)

type IpAllocator struct {
	IpPool map[string]bool
}
func NewIpAllocator() *IpAllocator {
	ipAllocator := &IpAllocator{
		IpPool: make(map[string]bool),
	}
	for i := 0; i < 254; i++ {
		ipAllocator.IpPool[BaseIp+fmt.Sprintf("%d", i)] = false
	}
	return ipAllocator
}
	
func (ipAllocator *IpAllocator) AllocateIp() string {
	for ip, used := range ipAllocator.IpPool {
		if !used {
			ipAllocator.IpPool[ip] = true
			return ip
		}
	}
	return ""
}

func (ipAllocator *IpAllocator) ReleaseIp(ip string) {
	ipAllocator.IpPool[ip] = false
}