package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/EmersonRabelo/report-processing-service/internal/api/perspective"
	contracts "github.com/EmersonRabelo/report-processing-service/internal/dto/report/contracts"
	"github.com/EmersonRabelo/report-processing-service/internal/entity"
	"github.com/EmersonRabelo/report-processing-service/internal/queue/producer"
	"github.com/google/uuid"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
}

type ConsumerReportService struct {
	repository  ReportRepository
	perspective perspective.PerspectiveAPIClient
	producer    producer.ReportAnalysisProducer
}

func NewConsumerReportService(repository ReportRepository, perspectiveAPIClient perspective.PerspectiveAPIClient, producer producer.ReportAnalysisProducer) *ConsumerReportService {
	return &ConsumerReportService{
		repository:  repository,
		perspective: perspectiveAPIClient,
		producer:    producer,
	}
}

func (crs *ConsumerReportService) Create(msg contracts.CreateReportMessage) error {

	if msg.Id == uuid.Nil {
		return fmt.Errorf("%w: Id is null", ErrInvalidMessage)
	}

	if msg.PostId == uuid.Nil {
		return fmt.Errorf("%w: PostId is null", ErrInvalidMessage)
	}

	if msg.ReporterId == uuid.Nil {
		return fmt.Errorf("%w: ReporterId is null", ErrInvalidMessage)
	}

	if msg.Body == "" {
		return fmt.Errorf("%w: Body is empty", ErrInvalidMessage)
	}

	createdAt := msg.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	report := entity.Report{
		Id:         msg.Id,
		PostId:     msg.PostId,
		ReporterId: msg.ReporterId,
		Status:     entity.StatusPending,
		CreatedAt:  createdAt,
	}

	crs.repository.InsertIfNotExists(&report)

	resp, err := crs.perspective.AnalyzePost(&msg.Body)

	if err != nil {
		return errors.New("Error processing content analysis")
	}

	fmt.Println(resp)

	return nil
}
