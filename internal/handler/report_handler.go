package handler

import (
	"encoding/json"
	"errors"

	contracts "github.com/EmersonRabelo/report-processing-service/internal/dto/report/contracts"
	"github.com/EmersonRabelo/report-processing-service/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrPermanent = errors.New("permanent error")

type ReportHandler struct {
	service *service.ConsumerReportService
}

func NewReportHandler(service *service.ConsumerReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (handler *ReportHandler) Handler(delivey amqp.Delivery) error {
	var msg contracts.CreateReportMessage
	if err := handler.jsonParse(&delivey.Body, &msg); err != nil {
		return ErrPermanent
	}

	if err := handler.service.Create(msg); err != nil {
		return err
	}

	return nil
}

func (handler *ReportHandler) jsonParse(body *[]byte, contract *contracts.CreateReportMessage) error {
	return json.Unmarshal(*body, &contract)
}
