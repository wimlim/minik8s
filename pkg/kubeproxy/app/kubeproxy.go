package kubeproxy

import (
	"encoding/json"
	"fmt"
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

func (kp *KubeProxy) handleServiceAdd(msg message.Message) {
}

func (kp *KubeProxy) handleServiceDelete(msg message.Message) {
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
