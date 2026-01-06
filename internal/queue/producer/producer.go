package producer

import (
	"encoding/json"
	"log"

	contract "github.com/EmersonRabelo/report-processing-service/internal/dto/report/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ReportAnalysisProducer struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

func NewReportAnalysisProducer(ch *amqp.Channel, exchange, routingKey string) *ReportAnalysisProducer {
	return &ReportAnalysisProducer{
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func (p *ReportAnalysisProducer) Publish(message *contract.ReportAnalysisResultMessage) error {
	body, err := json.Marshal(message)

	failOnError(err, "JSON parsing failed")

	return p.channel.Publish(
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
