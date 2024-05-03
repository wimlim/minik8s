package main

import(
	"minik8s/pkg/message"
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	// p := message.NewPublisher()
	// defer p.Close()
	// p.Publish(message.ScheduleQueue, []byte("hello"))
	s := message.NewSubscriber()
	defer s.Close()
	for{
		s.Subscribe(message.ScheduleQueue, func(d amqp.Delivery) {
			fmt.Println(string(d.Body))
		})
	
	}
}