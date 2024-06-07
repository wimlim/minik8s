package serverlessconfig

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	URL_HttpTrigger       = "/:namespace/:name"
	ServerDefaultListenIp = "0.0.0.0"
	ServerLocalIP         = "127.0.0.1"
)

var ServerMasterIP = ""
var ServerDefaultPort = 8081

type NodeConfig struct {
	ApiServerIP    string `yaml:"apiServerIP"`
	ApiServerPort  int    `yaml:"apiServerPort"`
	ServerlessIP   string `yaml:"serverlessIP"`
	ServerlessPort int    `yaml:"serverlessPort"`
}

func init() {
	// fmt.Println("serverConfig init")
	fd, err := os.Open("/config.yaml")
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

	var nodeConfig NodeConfig
	_ = yaml.Unmarshal(content, &nodeConfig)
	ServerMasterIP = nodeConfig.ApiServerIP
	ServerDefaultPort = nodeConfig.ApiServerPort
}

func GetMasterIP() string {
	return ServerMasterIP
}
func GetServerlessServerUrl() string {
	ServerURL := fmt.Sprintf("http://"+ServerMasterIP+":%d", ServerDefaultPort)
	return ServerURL
}
