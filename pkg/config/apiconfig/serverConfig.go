package apiconfig

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ServerDefaultListenIp = "0.0.0.0"
	HttpScheme            = "http://"
	ServerLocalIP         = "127.0.0.1"
)

var ServerMasterIP = ""
var ServerDefaultPort = 8080

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
func GetApiServerUrl() string {
	ServerURL := fmt.Sprintf(HttpScheme+ServerMasterIP+":%d", ServerDefaultPort)
	return ServerURL
}
