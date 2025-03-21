package events

import (
	"context"
	"fmt"
	"time"

	"github.com/notifier-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AdminEvent(name, email, phone string) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	helpers.RabbitError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	helpers.RabbitError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"admin",
		false,
		false,
		false,
		false,
		nil,
	)
	helpers.RabbitError(err, "Failed to declare a queue")
	body := fmt.Sprintf("%s, you have successfully created the admin with email: %s and phone: %s", name, email, phone)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	helpers.RabbitError(err, "Failed to publish a message")
}
