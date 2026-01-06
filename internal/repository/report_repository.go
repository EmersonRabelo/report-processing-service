package repository

import (
	"github.com/EmersonRabelo/report-processing-service/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
	Update(report *entity.Report) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) InsertIfNotExists(rep *entity.Report) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(rep).Error
}

func (r *reportRepository) Update(report *entity.Report) error {
	return r.db.Save(report).Error
}
