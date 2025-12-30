package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReportConsumer struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// Construtor
func NewReportConsumer(channel *amqp.Channel, exchange, routingKey string) *ReportConsumer {
	return &ReportConsumer{
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s, %s", msg, err)
	}
}

func (consumer *ReportConsumer) Consumer() {
	var exchangeType string = "topic"

	err := consumer.channel.ExchangeDeclare(
		consumer.exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	queue, err := consumer.channel.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare queue")

	err = consumer.channel.QueueBind(
		queue.Name,
		consumer.routingKey,
		consumer.exchange,
		false,
		nil,
	)

	failOnError(err, "Failed to bind queue")

	messages, err := consumer.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register consumer")

	log.Printf("Consumidor ativo! Binding: [%s]", consumer.routingKey)

	for message := range messages {
		log.Printf("Recebido: %s", message)
	}

}
