package kubeproxy

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/message"

	"github.com/streadway/amqp"
)

type KubeProxy struct {
	ipvsManager *IPVSManager
	subscriber  *message.Subscriber
}

func NewKubeProxy() *KubeProxy {
	ipvsManager := NewIPVSManager()
	if ipvsManager == nil {
		fmt.Println("Failed to create IPVS manager")
		return nil
	}
	subscriber := message.NewSubscriber()

	return &KubeProxy{
		ipvsManager: ipvsManager,
		subscriber:  subscriber,
	}
}

func podMatchesService(pod *apiobj.Pod, service *apiobj.Service) bool {
	labels := pod.MetaData.Labels
	for key, value := range service.Spec.Selector {
		if currentValue, ok := labels[key]; !ok || currentValue != value {
			return false
		}
	}
	return true
}

func (kp *KubeProxy) handleServiceAdd(msg message.Message) {
	var service apiobj.Service
	if err := json.Unmarshal([]byte(msg.Content), &service); err != nil {
		fmt.Println("Failed to unmarshal service:", err)
		return
	}

	pods, err := apirequest.GetAllPods()
	if err != nil {
		fmt.Println("Failed to get all pods:", err)
		return
	}

	var podIPs []string
	for _, pod := range pods {
		if podMatchesService(&pod, &service) {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}

	if len(podIPs) == 0 {
		fmt.Println("No pods match service selector")
		return
	}

	kp.ipvsManager.AddService(service.Spec, podIPs)
}

func (kp *KubeProxy) handleServiceDelete(msg message.Message) {
	var serviceid string
	if err := json.Unmarshal([]byte(msg.Content), &serviceid); err != nil {
		fmt.Println("Failed to unmarshal service ID:", err)
		return
	}
	kp.ipvsManager.DeleteService(serviceid)
}

func (kp *KubeProxy) handleServiceUpdate(msg message.Message) {
}

func (kp *KubeProxy) Run() {
	defer kp.subscriber.Close()
	defer kp.ipvsManager.Close()

	kp.subscriber.Subscribe(message.ServiceQueue, func(d amqp.Delivery) {
		var msg message.Message
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			fmt.Println("unmarshal message error:", err)
			return
		}

		switch msg.Type {
		case "Add":
			kp.handleServiceAdd(msg)
		case "Delete":
			kp.handleServiceDelete(msg)
		case "Update":
			kp.handleServiceUpdate(msg)
		}
	})
}
