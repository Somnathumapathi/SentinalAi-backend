package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Organization represents an organization in the system
type Organization struct {
	ID        string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name      string         `json:"name" gorm:"not null"`
	OwnerID   string         `json:"owner_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// AWSIntegration represents an AWS integration
type AWSIntegration struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrganizationID  string         `json:"organization_id" gorm:"not null"`
	AccessKeyID     string         `json:"access_key_id" gorm:"not null"`
	SecretAccessKey string         `json:"secret_access_key" gorm:"not null"`
	Region          string         `json:"region" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request models
type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAWSIntegrationRequest struct {
	OrganizationID  string `json:"organization_id" binding:"required"`
	AccessKeyID     string `json:"access_key_id" binding:"required"`
	SecretAccessKey string `json:"secret_access_key" binding:"required"`
	Region          string `json:"region" binding:"required"`
}
