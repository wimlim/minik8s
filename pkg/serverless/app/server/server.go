package server

import (
	"fmt"

	"encoding/json"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/config/serverlessconfig"
	"minik8s/pkg/serverless/app/autoscaler"
	"minik8s/tools/runner"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
)

type server struct {
	ip         string
	port       int
	router     *gin.Engine
	funcPodMap map[string][]string
}

func NewServer(ip string, port int) *server {
	return &server{
		ip:         ip,
		port:       port,
		router:     gin.Default(),
		funcPodMap: make(map[string][]string),
	}
}

func (s *server) Bind() {

	s.router.POST(serverlessconfig.URL_HttpTrigger, s.FunctionTrigger)
	s.router.GET(serverlessconfig.URL_HttpTrigger, s.FunctionCheck)

}

func (s *server) Run() {

	go runner.NewRunner().RunLoop(5*time.Second, 5*time.Second, s.FuncPodMapUpdateRoutine)
	go autoscaler.NewFuncScaler().Run()
	s.Bind()
	s.router.Run(fmt.Sprintf("%s:%d", s.ip, s.port))
}

func (s *server) FunctionTrigger(c *gin.Context) {
	func_namespace := c.Param("namespace")
	func_name := c.Param("name")

	key := func_namespace + "/" + func_name
	pod_ips, ok := s.funcPodMap[key]

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No available pod for function",
		})
		return
	}

	if len(pod_ips) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No available pod for function",
		})
		return
	}

	idx := rand.New(rand.NewSource(int64(len(pod_ips)))).Intn(len(pod_ips))
	pod_ip := pod_ips[idx]

	URL := fmt.Sprintf("http://%s:8080", pod_ip)
	body := c.Request.Body
	resp, err := http.Post(URL, "application/json", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Function trigger failed",
		})
		return
	}

	defer resp.Body.Close()
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Printf("decode res error\n")
		return
	}

	c.JSON(http.StatusOK, res["result"].(float64))
}

func (s *server) FunctionCheck(c *gin.Context) {
	func_namespace := c.Param("namespace")
	func_name := c.Param("name")

	key := func_namespace + "/" + func_name
	pod_ips, ok := s.funcPodMap[key]

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No available pod for function",
		})
		return
	}

	if len(pod_ips) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No available pod for function",
		})
		return
	}

}

func (s *server) FuncPodMapUpdateRoutine() {
	pods, err := apirequest.GetAllPods()
	if err != nil {
		fmt.Println("Get all pods failed")
		return
	}

	remoteFuncPodMap := make(map[string][]string)

	for _, pod := range pods {
		if pod.MetaData.Labels["func_uid"] != "" {
			key := pod.MetaData.Labels["func_namespace"] + "/" + pod.MetaData.Labels["func_name"]
			if pod.Status.PodIP != "" {
				remoteFuncPodMap[key] = append(remoteFuncPodMap[key], pod.Status.PodIP)
			}
		}
	}

	s.funcPodMap = remoteFuncPodMap
}
