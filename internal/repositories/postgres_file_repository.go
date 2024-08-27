package repositories

import (
	"test-go/internal/core/domain"

	"gorm.io/gorm"
)

type postgresFileRepository struct {
	db *gorm.DB
}

func NewPostgresFileRepository(db *gorm.DB) *postgresFileRepository {
	return &postgresFileRepository{db: db}
}

func (r *postgresFileRepository) Create(file *domain.FileAttachment) error {
	return r.db.Create(file).Error
}
