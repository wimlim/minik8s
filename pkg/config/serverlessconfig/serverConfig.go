package serverlessconfig

import (
	"fmt"
)

const (
	URL_HttpTrigger = "/:namespace/:name"

	ServerDefaultListenIp = "0.0.0.0"
	ServerDefaultPort     = 8081
	ServerLocalIP         = "127.0.0.1"
	ServerMasterIP        = "10.119.13.134"
)

func GetMasterIP() string {
	return ServerLocalIP
	// return ServerMasterIP
}
func GetServerlessServerUrl() string {

	serverless_ip := GetMasterIP()
	ServerURL := fmt.Sprintf("http://"+serverless_ip+":%d", ServerDefaultPort)
	return ServerURL
}
