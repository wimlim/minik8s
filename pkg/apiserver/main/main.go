package main

import(
	"fmt"
	"minik8s/pkg/apiserver/app/server"
	"minik8s/pkg/config"
)

func main() {
	fmt.Println("server start")
	s := server.NewServer(config.ServerDefaultListenIp, config.ServerDefaultPort)
	s = server.SetServer(s)
	s.Run()
}
