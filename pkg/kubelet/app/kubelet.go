package app

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	apiserverutil "minik8s/pkg/kubelet/app/apiserverUtil"
	"minik8s/pkg/kubelet/app/cache"
	podmanager "minik8s/pkg/kubelet/app/podManager"
	"minik8s/pkg/kubelet/app/status"
	"minik8s/pkg/message"
	prometheusutil "minik8s/pkg/prometheus/prometheusUtil"
	"os"

	"github.com/streadway/amqp"
)

type Kubelet struct {
	hostNode   *apiobj.Node
	podManager *podmanager.PodManager
	podCache   *cache.PodCache
}

type KubeletInterface interface {
	Run()
}

func NewKubelet() *Kubelet {
	newHostNode := getHostNodeConfig()
	newPodManager := podmanager.NewPodManager()
	newPodCache := cache.NewPodCache()
	return &Kubelet{
		hostNode:   newHostNode,
		podManager: newPodManager,
		podCache:   newPodCache,
	}
}

func (k *Kubelet) msgHandler(d amqp.Delivery) {
	fmt.Println(string(d.Body))
	var msg message.Message
	json.Unmarshal(d.Body, &msg)
	fmt.Println(msg.Name)
	var pod apiobj.Pod
	json.Unmarshal([]byte(msg.Content), &pod)
	if msg.Type == "Delete" {
		k.podManager.DeletePod(&pod)
		fmt.Println(pod.MetaData.Name)
	} else if msg.Type == "Add" {
		k.podManager.AddPod(&pod)
		fmt.Println(pod.MetaData.Name)
		apiserverutil.PodUpdate(&pod)
	}
}

func (k *Kubelet) listWatcher() {
	s := message.NewSubscriber()
	defer s.Close()

	hostname, _ := os.Hostname()
	que := fmt.Sprintf(message.PodQueue+"-%s", hostname)

	for {
		s.Subscribe(que, k.msgHandler)
	}
}

func (k *Kubelet) Run() {
	prometheusutil.StartPrometheusMetricsServer("10000")
	k.listWatcher()
	status.Run(k.podCache, k.hostNode)
}

func getHostNodeConfig() *apiobj.Node {
	allNodes, err := apiserverutil.GetAllNodes()
	if err != nil {
		fmt.Println("err:" + err.Error())
	}
	hostname, _ := os.Hostname()
	for _, node := range allNodes {
		if node.MetaData.Name == hostname {
			return &node
		}
	}
	return nil
}
