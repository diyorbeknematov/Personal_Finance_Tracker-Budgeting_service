package consumermq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer interface {
	ConsumeMessages(handle func(amqp.Delivery)) error
	Close() error
}

type rabbitmqConsumerImpl struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQConsumer(url string) (RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"budgeting", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &rabbitmqConsumerImpl{
		conn:      conn,
		channel:   ch,
		queueName: q.Name,
	}, nil
}


func (r *rabbitmqConsumerImpl) ConsumeMessages(handle func(amqp.Delivery)) error {
	msgs, err := r.channel.Consume(
		r.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handle(msg) 
		}
	}()

	return nil
}

func (r *rabbitmqConsumerImpl) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	if err := r.conn.Close(); err != nil {
		return err
	}
	return nil
}
