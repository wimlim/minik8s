package serviceconfig

import (
	"fmt"
)

const (
	BaseIp = "192.168.1."
)


var IpPool = make(map[string]bool)

func init() {
	for i := 0; i < 254; i++ {
		IpPool[BaseIp+fmt.Sprintf("%d", i)] = false
	}
}

func AllocateIp() string {
	
	for i := 0; i < 254; i++ {
		if _, ok := IpPool[BaseIp+fmt.Sprintf("%d", i)]; !ok {
			IpPool[BaseIp+fmt.Sprintf("%d", i)] = true
			return BaseIp + fmt.Sprintf("%d", i)
		}
	}
	return ""
}

func ReleaseIp(ip string) {
	delete(IpPool, ip)
}
