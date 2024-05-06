package server

import (
	"fmt"
	"minik8s/pkg/apiserver/app/handler"
	"minik8s/pkg/config/apiconfig"

	"github.com/gin-gonic/gin"
)

type server struct {
	ip     string
	port   int
	router *gin.Engine
}

func NewServer(ip string, port int) *server {
	return &server{
		ip:     ip,
		port:   port,
		router: gin.Default(),
	}
}

func SetServer(s *server) *server {
	s.Bind()
	return s
}

func (s *server) Bind() {
	//NODE
	s.router.GET(apiconfig.URL_AllNodes, handler.GetNodes)
	s.router.POST((apiconfig.URL_Node), handler.AddNode)
	s.router.DELETE((apiconfig.URL_Node), handler.DeleteNode)
	s.router.PUT((apiconfig.URL_Node), handler.UpdateNode)
	s.router.GET((apiconfig.URL_Node), handler.GetNode)
	s.router.GET((apiconfig.URL_NodeAllPods), handler.GetNodePods)
	s.router.GET((apiconfig.URL_NodeStatus), handler.GetNodeStatus)
	//POD
	s.router.GET((apiconfig.URL_GlobalPods), handler.GetGlobalPods)
	s.router.GET((apiconfig.URL_AllPods), handler.GetAllPods)
	s.router.POST((apiconfig.URL_Pod), handler.AddPod)
	s.router.DELETE((apiconfig.URL_Pod), handler.DeletePod)
	s.router.PUT((apiconfig.URL_Pod), handler.UpdatePod)
	s.router.GET((apiconfig.URL_Pod), handler.GetPod)
	s.router.GET((apiconfig.URL_PodStatus), handler.GetPodStatus)
	//SERVICE
	s.router.GET((apiconfig.URL_GlobalServices), handler.GetGlobalServices)
	s.router.GET((apiconfig.URL_AllServices), handler.GetAllServices)
	s.router.POST((apiconfig.URL_Service), handler.AddService)
	s.router.DELETE((apiconfig.URL_Service), handler.DeleteService)
	s.router.PUT((apiconfig.URL_Service), handler.UpdateService)
	s.router.GET((apiconfig.URL_Service), handler.GetService)
	s.router.GET((apiconfig.URL_ServiceStatus), handler.GetServiceStatus)

	fmt.Println("server bind success")
}

func (s *server) Run() {
	s.router.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}
