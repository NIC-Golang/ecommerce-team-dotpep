package consumer

import (
	"fmt"

	"github.com/notifier-service/internal/helpers"
	"github.com/notifier-service/internal/sender"
)

func AuthConsumer(email, queue string) {
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
		subject := fmt.Sprintf("Successful %s", q.Name)
		err := sender.SendEmail(email, subject, body)
		helpers.RabbitError(err, "Failed to send email")
	}
}
