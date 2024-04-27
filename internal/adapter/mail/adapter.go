package mail

import (
	"auth-service/internal/broker"
	"auth-service/internal/domain"
	"context"
	"encoding/json"

	rmq "github.com/rabbitmq/amqp091-go"
)

type (
	Adapter interface {
		SendMessage(message domain.EmailMessage) error
	}

	adapter struct {
		channel *rmq.Channel
		queue   rmq.Queue
	}
)

func NewAdapter(mailQueue string, client *broker.RabbitClient) (Adapter, error) {
	return adapter{
		channel: client.Chan,
		queue:   client.Queues[mailQueue],
	}, nil
}

func (a adapter) SendMessage(message domain.EmailMessage) error {
	messageStr, err := json.Marshal(message)

	if err != nil {
		return err
	}

	if err := a.channel.PublishWithContext(
		context.Background(),
		"",
		a.queue.Name,
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
