package apiconfig

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	ServerDefaultListenIp = "0.0.0.0"
	HttpScheme            = "http://"
	ServerLocalIP         = "127.0.0.1"
)

var ServerMasterIP = ""
var ServerDefaultPort = 8080

func init() {
	// fmt.Println("serverConfig init")
	fd, err := os.Open("/apiserver.txt")
	if err != nil {
		fmt.Println("open masterip.txt error")
		return
	}
	defer fd.Close()

	content, err := io.ReadAll(fd)
	if err != nil {
		fmt.Println("read masterip.txt error")
		return
	}

	ip_port := strings.Split(string(content), ":")
	ServerMasterIP = ip_port[0]
	_, err = fmt.Sscanf(ip_port[1], "%d", &ServerDefaultPort)
	if err != nil {
		fmt.Println("parse port error")
		return
	}
}

func GetMasterIP() string {
	return ServerMasterIP
}
func GetApiServerUrl() string {
	ServerURL := fmt.Sprintf(HttpScheme+ServerMasterIP+":%d", ServerDefaultPort)
	return ServerURL
}
