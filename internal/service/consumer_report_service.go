package service

import (
	"errors"
	"fmt"
	"time"

	contracts "github.com/EmersonRabelo/report-processing-service/internal/dto/report/contracts"
	"github.com/EmersonRabelo/report-processing-service/internal/entity"
	"github.com/google/uuid"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
}

type ConsumerReportService struct {
	repository ReportRepository
}

func NewConsumerReportService(repository ReportRepository) *ConsumerReportService {
	return &ConsumerReportService{
		repository: repository,
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

	return crs.repository.InsertIfNotExists(&report)
}
