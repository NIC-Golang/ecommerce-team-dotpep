package helpers

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitConnecter(queueName string) (*amqp.Connection, *amqp.Channel, amqp.Queue, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		queueName, // имя очереди
		false,     // очередь не устойчива
		false,     // не удаляется, если не используется
		false,     // не эксклюзивная
		false,     // нет ожидания
		nil,       // дополнительные параметры
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, amqp.Queue{}, err
	}

	return conn, ch, q, nil
}
func RabbitError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
