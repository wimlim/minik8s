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
	monitormanager "minik8s/pkg/prometheus/monitorManager"

	"github.com/streadway/amqp"
)

type Kubelet struct {
	podManager *podmanager.PodManager
	podCache   *cache.PodCache
}

type KubeletInterface interface {
	Run()
}

func NewKubelet() *Kubelet {
	newPodManager := podmanager.NewPodManager()
	newPodCache := cache.NewPodCache()
	return &Kubelet{
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
		monitormanager.RemovePodMonitor(&pod)
		fmt.Println(pod.MetaData.Name)
	} else if msg.Type == "Add" {
		k.podManager.AddPod(&pod)
		monitormanager.AddPodMonitor(&pod)
		fmt.Println(pod.MetaData.Name)
		apiserverutil.PodUpdate(&pod)
	}
}

func (k *Kubelet) listWatcher() {
	s := message.NewSubscriber()
	defer s.Close()
	for {
		s.Subscribe(message.PodQueue, k.msgHandler)
	}
}

func (k *Kubelet) Run() {
	status.Run(k.podCache)
	k.listWatcher()
}
