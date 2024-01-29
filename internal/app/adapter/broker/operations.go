package broker

import (
	"context"
	"encoding/json"

	rmq "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	Connection *rmq.Connection
	Channel    *rmq.Channel
}

func (b Broker) SendMessage(queueName string, message interface{}) error {
	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err := b.Channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		rmq.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageStr),
		},
	); err != nil {
		return err
	}

	return nil
}
