package message

import (
	"github.com/streadway/amqp"
)

type Subscriber struct {
	connct *amqp.Connection
	url    string
}

func NewSubscriber() *Subscriber {
	url := RabbitMQURL()
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	return &Subscriber{
		connct: conn,
		url:    url,
	}
}
func (s *Subscriber) Close() {
	s.connct.Close()
}
func (s *Subscriber) Subscribe(routingKey string, callback func(amqp.Delivery)) error {

	ch, err := s.connct.Channel()
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

	msgs, err := ch.Consume(routingKey, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for d := range msgs {
		callback(d)
	}
	return nil
}
