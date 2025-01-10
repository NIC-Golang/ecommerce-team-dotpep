package events

import (
	"context"
	"fmt"
	"time"

	"github.com/notifier-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SignUpEvent(username, email string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	helpers.RabbitError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	helpers.RabbitError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"signup",
		false,
		false,
		false,
		false,
		nil,
	)
	helpers.RabbitError(err, "Failed to declare a queue")
	body := fmt.Sprintf("%s, you have successfully signed up with email %s", username, email)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	helpers.RabbitError(err, "Failed to publish a message")
}

func LoginEvent(username string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	helpers.RabbitError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	helpers.RabbitError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"login",
		false,
		false,
		false,
		false,
		nil,
	)
	helpers.RabbitError(err, "Failed to declare a queue")
	body := fmt.Sprintf("%s, you have successfully logged in", username)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body)})
	helpers.RabbitError(err, "Failed to publish a message")
}
