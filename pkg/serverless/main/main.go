package main

import (
	"minik8s/pkg/serverless/app/server"
	"minik8s/pkg/config/serverlessconfig"
)

func main() {
	s := server.NewServer(serverlessconfig.ServerDefaultListenIp, serverlessconfig.ServerDefaultPort)
	s.Run()
}