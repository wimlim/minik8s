package main

import(
	"fmt"
	"minik8s/pkg/apiserver/app/server"
	"minik8s/pkg/config/apiconfig"
)

func main() {
	fmt.Println("server start")
	s := server.NewServer(apiconfig.ServerDefaultListenIp, 
						apiconfig.ServerDefaultPort)
	s = server.SetServer(s)
	s.Run()
}
