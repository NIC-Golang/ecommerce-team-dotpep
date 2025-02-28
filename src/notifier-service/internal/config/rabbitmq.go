package config

import (
	"log"

	"github.com/notifier-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQ() {
	_, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	helpers.RabbitError(err, "Failed to connect to RabbitMQ")
	log.Printf("Connected to RabbitMQ")
}
