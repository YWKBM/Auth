package queue

import (
	"auth/config"
	"auth/queue/messages"
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]string
}

func AddQueue(rabbitMqConfig config.RabbitMQConfig) (*Queue, error) {
	conn, err := amqp.Dial(rabbitMqConfig.RABBIT_URL)
	if err != nil {
		return nil, err
	}

	chanel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Queue{
		conn:    conn,
		channel: chanel,
		queues:  make(map[string]string),
	}, nil
}

func (q *Queue) CreateQueue(queueName string) error {
	_, err := q.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q.queues[queueName] = queueName
	log.Printf("Queue %s has been added.", queueName)
	return nil
}

func (q *Queue) AddConsumer(queueName string, consumeFunc func([]byte)) error {
	if _, exists := q.queues[queueName]; !exists {
		return amqp.ErrClosed // Очередь не найдена
	}

	msgs, err := q.channel.Consume(
		queueName, // имя очереди
		"",        // имя потребителя
		true,      // авто-подтверждение
		false,     // эксклюзивный
		false,     // локальный
		false,     // без ожидания
		nil,       // аргументы
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			consumeFunc(d.Body)
		}
	}()

	log.Printf("Consumer added to queue %s.", queueName)
	return nil
}

func (q *Queue) SendMessage(message messages.Message) error {
	if _, exists := q.queues[message.RoutingKey]; !exists {
		return amqp.ErrClosed
	}

	err := q.channel.Publish(
		"",
		message.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message.Body,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Message sent to queue %s: %s", message.RoutingKey, message)
	return nil
}

func (q *Queue) Close() {
	if q.channel != nil {
		_ = q.channel.Close()
	}
	if q.conn != nil {
		_ = q.conn.Close()
	}
	log.Println("Queue connection closed.")
}
