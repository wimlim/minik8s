package workflow

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/message"

	"github.com/streadway/amqp"
)

func StartWorkflow(w *apiobj.Workflow) {
	fmt.Println("Start workflow")
}

func Run() {
	sub := message.NewSubscriber()
	defer sub.Close()
	sub.Subscribe(message.WorkflowQueue, func(d amqp.Delivery) {
		fmt.Println("Receive workflow message")
		var message message.Message
		json.Unmarshal([]byte(d.Body), &message)

		var workflow apiobj.Workflow
		json.Unmarshal([]byte(message.Content), &workflow)

		StartWorkflow(&workflow)
	})
}
