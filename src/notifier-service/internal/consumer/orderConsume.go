package consumer

import (
	"github.com/notifier-service/internal/helpers"
	"github.com/notifier-service/internal/sender"
)

func OrderConsumer(email, queue string) {
	conn, ch, q, err := helpers.RabbitConnecter(queue)
	helpers.RabbitError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer ch.Close()
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	helpers.RabbitError(err, "Failed to register a consumer")
	for msg := range msgs {
		body := string(msg.Body)
		subject := "Successful creating of order!"
		err := sender.SendEmail(email, subject, body)
		helpers.RabbitError(err, "Failed to send email")
	}
}
