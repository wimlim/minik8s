package message

import (
	"github.com/streadway/amqp"
)

type Publisher struct {
	connct *amqp.Connection
	url    string
}

func NewPublisher() *Publisher {
	url := RabbitMQURL()
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return &Publisher{
		connct: conn,
		url:    url,
	}
}

func (p *Publisher) Close() {
	p.connct.Close()
}

func (p *Publisher) Publish(routingKey string, msg []byte) error {
	ch, err := p.connct.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(routingKey, true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(routingKey, routingKey, Queue2Exchange[routingKey], false, nil)
	if err != nil {
		return err
	}
	err = ch.Publish(
		Queue2Exchange[routingKey],
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	if err != nil {
		return err
	}
	return nil
}
