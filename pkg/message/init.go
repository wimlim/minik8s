package message

import (
	"github.com/streadway/amqp"
)

func init() {
	url := RabbitMQURL()
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"minik8s",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

}
