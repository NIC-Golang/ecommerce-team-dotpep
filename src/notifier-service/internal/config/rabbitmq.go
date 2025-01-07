package config

import (
	"context"
	"log"
	"time"

	"github.com/notifier-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	helpers.RabbitError(err, "Error with connecting to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	helpers.RabbitError(err, "Error with opening the channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Hello, Rabbit!",
		false,
		false,
		false,
		false,
		nil,
	)

	helpers.RabbitError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	body := "Hello, RabbitMQ!"
	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	helpers.RabbitError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
