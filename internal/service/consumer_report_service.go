package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/EmersonRabelo/report-processing-service/internal/api/perspective"
	contracts "github.com/EmersonRabelo/report-processing-service/internal/dto/report/contracts"
	"github.com/EmersonRabelo/report-processing-service/internal/entity"
	"github.com/EmersonRabelo/report-processing-service/internal/queue/producer"
	"github.com/EmersonRabelo/report-processing-service/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
}

type ConsumerReportService struct {
	repository  repository.ReportRepository
	perspective perspective.PerspectiveAPIClient
	producer    producer.ReportAnalysisProducer
}

func NewConsumerReportService(repository repository.ReportRepository, perspectiveAPIClient perspective.PerspectiveAPIClient, producer producer.ReportAnalysisProducer) *ConsumerReportService {
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

	if err := crs.repository.InsertIfNotExists(&report); err != nil {
		return errors.New("Error persist the report")
	}

	resp, err := crs.perspective.AnalyzePost(&msg.Body)

	if err != nil {
		return errors.New("Error processing content analysis")
	}

	for attribute, score := range map[string]float64{
		"TOXICITY":        resp.AttributeScores.Toxicity.SummaryScore.Value,
		"SEVERE_TOXICITY": resp.AttributeScores.SevereToxicity.SummaryScore.Value,
		"IDENTITY_ATTACK": resp.AttributeScores.IdentityAttack.SummaryScore.Value,
		"INSULT":          resp.AttributeScores.Insult.SummaryScore.Value,
		"PROFANITY":       resp.AttributeScores.Profanity.SummaryScore.Value,
		"THREAT":          resp.AttributeScores.Threat.SummaryScore.Value,
	} {
		switch attribute {
		case "TOXICITY":
			report.PerspectiveToxicity = ptr(score)
		case "SEVERE_TOXICITY":
			report.PerspectiveSevereToxicity = ptr(score)
		case "IDENTITY_ATTACK":
			report.PerspectiveIdentityAttack = ptr(score)
		case "INSULT":
			report.PerspectiveInsult = ptr(score)
		case "PROFANITY":
			report.PerspectiveProfanity = ptr(score)
		case "THREAT":
			report.PerspectiveThreat = ptr(score)

		}
	}

	if len(resp.DetectedLanguages) > 0 {
		report.PerspectiveLanguage = &resp.DetectedLanguages[0]
	}

	now := time.Now()

	report.UpdatedAt = now
	report.PerspectiveResponseAt = &now
	report.Status = entity.StatusDone

	if err := crs.repository.Update(&report); err != nil {
		return errors.New("Error persist content analysis")
	}

	reportAnalysisMessage := ToReportAnalysisResultMessage(report)

	if err := crs.producer.Publish(&reportAnalysisMessage); err != nil {
		return errors.New("Error publish analysis menssage")
	}

	fmt.Println("Tudo Certo!!!")
	fmt.Printf("Message: %+v\n", reportAnalysisMessage)

	return nil
}

func ToReportAnalysisResultMessage(r entity.Report) contracts.ReportAnalysisResultMessage {
	return contracts.ReportAnalysisResultMessage{
		ReportId:       r.Id,
		Toxicity:       r.PerspectiveToxicity,
		SevereToxicity: r.PerspectiveSevereToxicity,
		IdentityAttack: r.PerspectiveIdentityAttack,
		Insult:         r.PerspectiveInsult,
		Profanity:      r.PerspectiveProfanity,
		Threat:         r.PerspectiveThreat,
		Language:       r.PerspectiveLanguage,
		AnalyzedAt:     r.PerspectiveResponseAt,
	}
}

// Fora idiomática, posso não receber os valores corretamente. Então checo e preencho os valores.
func ptr[T any](v T) *T { return &v }
