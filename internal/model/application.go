package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	ApplicationStatusActive    ApplicationStatus = "active"
	ApplicationStatusInactive  ApplicationStatus = "inactive"
	ApplicationStatusSuspended ApplicationStatus = "suspended"
)

type Application struct {
	ID          uuid.UUID         `gorm:"primaryKey"`
	Name        string            `gorm:"not null"`
	Description string            `gorm:"not null"`
	APIKey      string            `json:"api_key" gorm:"uniqueIndex;not null"`
	Status      ApplicationStatus `json:"status" gorm:"default:'active'"`
	UserID      uuid.UUID         `json:"user_id" gorm:"not null"`
	User        User              `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `json:"-" gorm:"index"` // Soft delete
}

// TableName specifies the table name for the Application model
func (Application) TableName() string {
	return "applications"
}

// BeforeCreate is a GORM hook that runs before creating a record
func (a *Application) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID for the application
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	// Generate API key
	if a.APIKey == "" {
		a.APIKey = generateAPIKey()
	}

	return nil
}

// generateAPIKey generates a secure API key
func generateAPIKey() string {
	return "app_" + uuid.New().String()
}
