package apiconfig

import "fmt"

const (
	ServerDefaultListenIp = "0.0.0.0"
	ServerDefaultPort     = 8080
	HttpScheme            = "http://"
	ServerLocalIP         = "127.0.0.1"
	ServerMasterIP        = "10.119.13.134"
)

func GetMasterIP() string {
	return ServerLocalIP
	// return ServerMasterIP
}
func GetApiServerUrl() string {

	apiserver_ip := GetMasterIP()
	ServerURL := fmt.Sprintf(HttpScheme+apiserver_ip+":%d", ServerDefaultPort)
	return ServerURL
}
