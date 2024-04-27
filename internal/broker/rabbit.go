package broker

import (
	"fmt"

	rmq "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Queues map[string]rmq.Queue
	Conn   *rmq.Connection
	Chan   *rmq.Channel
}

// Добавить коннекшн через TLS как с постгрой
func NewRabbitClient(URL string, queues ...string) (*RabbitClient, error) {
	conn, err := rmq.Dial(URL)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to rabbit %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error while opening the channel %w", err)
	}

	qs := make(map[string]rmq.Queue, len(queues))
	for _, queue := range queues {
		q, err := ch.QueueDeclare(
			queue,
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return nil, fmt.Errorf("error while declaring the queues %w", err)
		}

		qs[queue] = q
	}

	return &RabbitClient{
		Queues: qs,
		Conn:   conn,
		Chan:   ch,
	}, nil
}
