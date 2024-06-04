package serverlessconfig

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	URL_HttpTrigger       = "/:namespace/:name"
	ServerDefaultListenIp = "0.0.0.0"
	ServerLocalIP         = "127.0.0.1"
)

var ServerMasterIP = ""
var ServerDefaultPort = "8081"

func init() {
	// fmt.Println("serverConfig init")
	fd, err := os.Open("/serverless.txt")
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
func GetServerlessServerUrl() string {
	ServerURL := fmt.Sprintf("http://"+ServerMasterIP+":%d", ServerDefaultPort)
	return ServerURL
}
