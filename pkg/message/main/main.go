package main

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
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
	fmt.Println(pod.MetaData.Name)

}
func main() {
	// p := message.NewPublisher()
	// defer p.Close()
	// p.Publish(message.ScheduleQueue, []byte("hello"))

	//This is a example
	s := message.NewSubscriber()
	defer s.Close()
	for {
		s.Subscribe(message.ScheduleQueue, msgHandler)
	}
}
