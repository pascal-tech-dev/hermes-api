package repository

import (
	"context"
	"hermes-api/internal/model"

	"gorm.io/gorm"
)

// ApplicationRepository defines the interface for application data operations
type ApplicationRepository interface {

	// Basic CRUD operations
	BaseRepository[model.Application]

	// Query operations
	GetByAPIKey(ctx context.Context, apiKey string) (*model.Application, error)
}

type applicationRepository struct {
	BaseRepository[model.Application]
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{
		BaseRepository: NewBaseRepository[model.Application](db),
		db:             db,
	}
}

func (r *applicationRepository) GetByAPIKey(ctx context.Context, apiKey string) (*model.Application, error) {
	var application model.Application
	err := r.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&application).Error
	if err != nil {
		return nil, err
	}
	return &application, nil
}
