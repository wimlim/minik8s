package serverlessconfig

import (
	"fmt"
)

const (
	RegistryServerBindIp   = "0.0.0.0"
	RegistryServerPort = 5000

	RegistryImage         = "registry:2.7.1"
	RegistryContainerName = "minik8s-registry"
)

func GetRegistryServerUrl() string {

	RegistryServerIp := GetMasterIP()
	return fmt.Sprintf("%s:%d", RegistryServerIp, RegistryServerPort)
}