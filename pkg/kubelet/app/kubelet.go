package app

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/kubelet/app/runtime"
	"minik8s/pkg/message"

	"github.com/streadway/amqp"
)

func msgHandler(d amqp.Delivery) {
	fmt.Println(string(d.Body))
	var msg message.Message
	json.Unmarshal(d.Body, &msg)
	fmt.Println(msg.Name)

	var pod apiobj.Pod
	json.Unmarshal([]byte(msg.Content), &pod)
	runtime.RestartPauseContainer(&pod)
	fmt.Println(pod.MetaData.Name)
}

func Run() {
	s := message.NewSubscriber()
	defer s.Close()
	for {
		s.Subscribe(message.PodQueue, msgHandler)
	}
}
