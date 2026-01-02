package repository

import (
	"github.com/EmersonRabelo/report-processing-service/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) InsertIfNotExists(rep *entity.Report) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(rep).Error
}
