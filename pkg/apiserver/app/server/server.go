package server

import (
	"fmt"
	"minik8s/pkg/apiserver/app/handler"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/config/serviceconfig"
	"os/exec"
	"minik8s/pkg/etcd"
	"mink8s/tools/weave"

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

func (s *server) StartDnsNginx() {
	cmd := exec.Command("docker", "run", "-d", "--name", "my-nginx-container", "-p", "80:80", "-v", "/root/minik8s/pkg/nginx:/etc/nginx/conf.d", "nginx")
    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error executing command:", err)
        return
    }
	ip, err := weave.WeaveAttach(string(output))
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	etcd.EtcdKV.Put(etcd.PATH_EtcdDnsNginxIP, []byte(ip))
	fmt.Println(ip)
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
	//REPLICASET
	s.router.GET((apiconfig.URL_GlobalReplicaSets), handler.GetGlobalReplicaSets)
	s.router.GET((apiconfig.URL_AllReplicaSets), handler.GetAllReplicaSets)
	s.router.POST((apiconfig.URL_ReplicaSet), handler.AddReplicaSet)
	s.router.DELETE((apiconfig.URL_ReplicaSet), handler.DeleteReplicaSet)
	s.router.PUT((apiconfig.URL_ReplicaSet), handler.UpdateReplicaSet)
	s.router.GET((apiconfig.URL_ReplicaSet), handler.GetReplicaSet)
	s.router.GET((apiconfig.URL_ReplicaSetStatus), handler.GetReplicaSetStatus)
	//HPA
	s.router.GET((apiconfig.URL_GlobalHpas), handler.GetGlobalHpas)
	s.router.GET((apiconfig.URL_AllHpas), handler.GetAllHpas)
	s.router.POST((apiconfig.URL_Hpa), handler.AddHpa)
	s.router.DELETE((apiconfig.URL_Hpa), handler.DeleteHpa)
	s.router.PUT((apiconfig.URL_Hpa), handler.UpdateHpa)
	s.router.GET((apiconfig.URL_Hpa), handler.GetHpa)
	s.router.GET((apiconfig.URL_HpaStatus), handler.GetHpaStatus)
	//DNS
	s.router.GET((apiconfig.URL_GlobalDns), handler.GetGlobalDns)
	s.router.GET((apiconfig.URL_AllDns), handler.GetAllDns)
	s.router.POST((apiconfig.URL_Dns), handler.AddDns)
	s.router.DELETE((apiconfig.URL_Dns), handler.DeleteDns)
	s.router.PUT((apiconfig.URL_Dns), handler.UpdateDns)
	s.router.GET((apiconfig.URL_Dns), handler.GetDns)
	s.router.GET((apiconfig.URL_DnsStatus), handler.GetDnsStatus)

	fmt.Println("server bind success")
}

func (s *server) Run() {
	serviceconfig.NewIpAllocator()
	s.StartDnsNginx()
	s.router.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}
