package service

import (
	"context"
	"hermes-api/internal/model"
	"hermes-api/internal/repository"

	"github.com/google/uuid"
)

// ApplicationService defines the interface for application business logic
type ApplicationService interface {
	CreateApplication(ctx context.Context, application *model.Application) error
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*model.Application, error)
	UpdateApplication(ctx context.Context, application *model.Application) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error
}

// applicationService implements ApplicationService
type applicationService struct {
	applicationRepo repository.ApplicationRepository
}

// NewApplicationService creates a new application service
func NewApplicationService(applicationRepo repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		applicationRepo: applicationRepo,
	}
}

// CreateApplication implements ApplicationService.
func (a *applicationService) CreateApplication(ctx context.Context, application *model.Application) error {

	panic("unimplemented")
}

// DeleteApplication implements ApplicationService.
func (a *applicationService) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetApplicationByID implements ApplicationService.
func (a *applicationService) GetApplicationByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	panic("unimplemented")
}

// UpdateApplication implements ApplicationService.
func (a *applicationService) UpdateApplication(ctx context.Context, application *model.Application) error {
	panic("unimplemented")
}
