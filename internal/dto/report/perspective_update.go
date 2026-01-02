package report

import (
	"time"

	"github.com/EmersonRabelo/report-processing-service/internal/entity"
)

type PerspectiveUpdate struct {
	Toxicity     *float64
	Insult       *float64
	Profanity    *float64
	Threat       *float64
	IdentityHate *float64
	Language     *string
	ResponseAt   *time.Time
	Status       entity.ProcessStatus
}
