package serviceconfig

import (
	"fmt"
	"encoding/json"
	"minik8s/pkg/etcd"
)

const (
	BaseIp = "192.168.1."
)

type IpAllocator struct {
	IpPool map[string]bool `json:"ipPool"`
}
func NewIpAllocator(){
	ipAllocator := &IpAllocator{
		IpPool: make(map[string]bool),
	}
	for i := 0; i < 254; i++ {
		ipAllocator.IpPool[BaseIp+fmt.Sprintf("%d", i)] = false
	}
	ipAllocatorJson, err := json.Marshal(ipAllocator)
	if err != nil {
		fmt.Println("json marshal error")
	}
	etcd.EtcdKV.Put(etcd.PATH_EtcdServieIps, ipAllocatorJson)
}
	
func  AllocateIp() string {
	ipAllocator := &IpAllocator{}
	key := etcd.PATH_EtcdServieIps
	ipAllocatorJson, err := etcd.EtcdKV.Get(key)
	if err != nil {
		fmt.Println("get ip pool error")
	}
	err = json.Unmarshal([]byte(ipAllocatorJson), ipAllocator)
	if err != nil {
		fmt.Println("json unmarshal error")
	}
	for ip, used := range ipAllocator.IpPool {
		if !used {
			ipAllocator.IpPool[ip] = true
			ipAllocatorJson, err = json.Marshal(ipAllocator)
			if err != nil {
				fmt.Println("json marshal error")
			}
			etcd.EtcdKV.Put(etcd.PATH_EtcdServieIps, ipAllocatorJson)
			return ip
		}
	}
	return ""
}

func ReleaseIp(ip string) {
	ipAllocator := &IpAllocator{}
	key := etcd.PATH_EtcdServieIps
	ipAllocatorJson, err := etcd.EtcdKV.Get(key)
	if err != nil {
		fmt.Println("get ip pool error")
	}
	err = json.Unmarshal([]byte(ipAllocatorJson), ipAllocator)
	if err != nil {
		fmt.Println("json unmarshal error")
	}
	ipAllocator.IpPool[ip] = false
	ipAllocatorJson, err = json.Marshal(ipAllocator)
	if err != nil {
		fmt.Println("json marshal error")
	}
	etcd.EtcdKV.Put(etcd.PATH_EtcdServieIps, ipAllocatorJson)
}