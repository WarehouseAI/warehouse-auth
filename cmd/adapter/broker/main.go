package broker

import (
	"auth-service/configs"
	"auth-service/internal/app/adapter/broker"
	"fmt"

	rmq "github.com/rabbitmq/amqp091-go"
)

func NewBroker(config configs.Config) *broker.Broker {
	conn, err := rmq.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.Rabbit.Username,
		config.Rabbit.Password,
		config.Rabbit.Host,
		config.Rabbit.Port,
	))

	if err != nil {
		panic(fmt.Sprintf("Unable to open connect to RabbitMQ; %s", err))
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(fmt.Sprintf("Unable to open channel; %s", err))
	}

	for _, queue := range config.Rabbit.Queues {
		_, err := ch.QueueDeclare(
			queue,
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			panic(fmt.Sprintf("Unable to create queue in the channel; %s", err))
		}
	}

	return &broker.Broker{
		Connection: conn,
		Channel:    ch,
	}
}
