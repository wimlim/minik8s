package server

import (
	"fmt"
	"minik8s/pkg/apiserver/app/handler"
	"minik8s/pkg/config/apiconfig"
	"minik8s/pkg/config/serviceconfig"
	"minik8s/pkg/etcd"
	"minik8s/tools/weave"
	"os/exec"

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
	cmd := exec.Command("docker", "run", "-d", "--privileged", "--name", "my-nginx-container", "-p", "80:80", "-v", "/root/minik8s/pkg/nginx:/etc/nginx/conf.d", "nginx")
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

	updateCmd := exec.Command("docker", "exec", "my-nginx-container", "apt", "update")
	_, err = updateCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	installCmd := exec.Command("docker", "exec", "my-nginx-container", "apt", "install", "-y", "ipvsadm")
	_, err = installCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
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
	s.router.PUT((apiconfig.URL_PodStatus), handler.UpdatePodStatus)
	//SERVICE
	s.router.GET((apiconfig.URL_GlobalServices), handler.GetGlobalServices)
	s.router.GET((apiconfig.URL_AllServices), handler.GetAllServices)
	s.router.POST((apiconfig.URL_Service), handler.AddService)
	s.router.DELETE((apiconfig.URL_Service), handler.DeleteService)
	s.router.PUT((apiconfig.URL_Service), handler.UpdateService)
	s.router.GET((apiconfig.URL_Service), handler.GetService)
	s.router.GET((apiconfig.URL_ServiceStatus), handler.GetServiceStatus)
	s.router.PUT((apiconfig.URL_ServiceStatus), handler.UpdateServiceStatus)
	//REPLICASET
	s.router.GET((apiconfig.URL_GlobalReplicaSets), handler.GetGlobalReplicaSets)
	s.router.GET((apiconfig.URL_AllReplicaSets), handler.GetAllReplicaSets)
	s.router.POST((apiconfig.URL_ReplicaSet), handler.AddReplicaSet)
	s.router.DELETE((apiconfig.URL_ReplicaSet), handler.DeleteReplicaSet)
	s.router.PUT((apiconfig.URL_ReplicaSet), handler.UpdateReplicaSet)
	s.router.GET((apiconfig.URL_ReplicaSet), handler.GetReplicaSet)
	s.router.GET((apiconfig.URL_ReplicaSetStatus), handler.GetReplicaSetStatus)
	s.router.PUT((apiconfig.URL_ReplicaSetStatus), handler.UpdateReplicaSetStatus)
	//HPA
	s.router.GET((apiconfig.URL_GlobalHpas), handler.GetGlobalHpas)
	s.router.GET((apiconfig.URL_AllHpas), handler.GetAllHpas)
	s.router.POST((apiconfig.URL_Hpa), handler.AddHpa)
	s.router.DELETE((apiconfig.URL_Hpa), handler.DeleteHpa)
	s.router.PUT((apiconfig.URL_Hpa), handler.UpdateHpa)
	s.router.GET((apiconfig.URL_Hpa), handler.GetHpa)
	s.router.GET((apiconfig.URL_HpaStatus), handler.GetHpaStatus)
	s.router.PUT((apiconfig.URL_HpaStatus), handler.UpdateHpaStatus)
	//DNS
	s.router.GET((apiconfig.URL_GlobalDns), handler.GetGlobalDns)
	s.router.GET((apiconfig.URL_AllDns), handler.GetAllDns)
	s.router.POST((apiconfig.URL_Dns), handler.AddDns)
	s.router.DELETE((apiconfig.URL_Dns), handler.DeleteDns)
	s.router.PUT((apiconfig.URL_Dns), handler.UpdateDns)
	s.router.GET((apiconfig.URL_Dns), handler.GetDns)
	s.router.GET((apiconfig.URL_DnsStatus), handler.GetDnsStatus)
	s.router.PUT((apiconfig.URL_DnsStatus), handler.UpdateDnsStatus)
	//FUNCTION
	s.router.GET((apiconfig.URL_GlobalFunctions), handler.GetGlobalFunctions)
	s.router.GET((apiconfig.URL_AllFunctions), handler.GetAllFunctions)
	s.router.POST((apiconfig.URL_Function), handler.AddFunction)
	s.router.DELETE((apiconfig.URL_Function), handler.DeleteFunction)
	s.router.PUT((apiconfig.URL_Function), handler.UpdateFunction)
	s.router.GET((apiconfig.URL_Function), handler.GetFunction)
	//WORKFLOW
	s.router.GET((apiconfig.URL_GlobalWorkflows), handler.GetGlobalWorkflows)
	s.router.GET((apiconfig.URL_AllWorkflows), handler.GetAllWorkflows)
	s.router.POST((apiconfig.URL_Workflow), handler.AddWorkflow)
	s.router.DELETE((apiconfig.URL_Workflow), handler.DeleteWorkflow)
	s.router.PUT((apiconfig.URL_Workflow), handler.UpdateWorkflow)
	s.router.GET((apiconfig.URL_Workflow), handler.GetWorkflow)
	s.router.PUT((apiconfig.URL_WorkflowStatus), handler.UpdateWorkflowStatus)
	//PV
	s.router.GET((apiconfig.URL_GlobalPVs), handler.GetGlobalPVs)
	s.router.GET((apiconfig.URL_AllPVs), handler.GetAllPVs)
	s.router.POST((apiconfig.URL_PV), handler.AddPV)
	s.router.DELETE((apiconfig.URL_PV), handler.DeletePV)
	s.router.PUT((apiconfig.URL_PV), handler.UpdatePV)
	s.router.GET((apiconfig.URL_PV), handler.GetPV)
	//PVC
	s.router.GET((apiconfig.URL_GlobalPVCs), handler.GetGlobalPVCs)
	s.router.GET((apiconfig.URL_AllPVCs), handler.GetAllPVCs)
	s.router.POST((apiconfig.URL_PVC), handler.AddPVC)
	s.router.DELETE((apiconfig.URL_PVC), handler.DeletePVC)
	s.router.PUT((apiconfig.URL_PVC), handler.UpdatePVC)
	s.router.GET((apiconfig.URL_PVC), handler.GetPVC)

	fmt.Println("server bind success")
}

func (s *server) Run() {
	serviceconfig.NewIpAllocator()
	s.StartDnsNginx()
	s.router.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}
