package models

import (
	"time"

	"gorm.io/gorm"
)

// Subscription represents a subscription in the system
type Subscription struct {
	ID                     string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrganizationID         string         `json:"organization_id" gorm:"not null"`
	Plan                   string         `json:"plan" gorm:"not null;check:plan IN ('free', 'pro', 'enterprise')"`
	Status                 string         `json:"status" gorm:"not null;check:status IN ('pending', 'active', 'inactive', 'suspended', 'cancelled')"`
	RazorpaySubscriptionID string         `json:"razorpay_subscription_id"`
	StartDate              time.Time      `json:"start_date" gorm:"not null"`
	EndDate                time.Time      `json:"end_date" gorm:"not null"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `json:"-" gorm:"index"`
}

// PlanLimits defines the limits for each subscription plan
type PlanLimits struct {
	GitHubRepos int `json:"github_repos"`
	AWSAccounts int `json:"aws_accounts"`
	TeamMembers int `json:"team_members"`
	CustomRules int `json:"custom_rules"`
	APIRequests int `json:"api_requests"`
}

// PlanLimitsMap defines the limits for each plan type
var PlanLimitsMap = map[string]PlanLimits{
	"free": {
		GitHubRepos: 1,
		AWSAccounts: 1,
		TeamMembers: 2,
		CustomRules: 5,
		APIRequests: 1000,
	},
	"pro": {
		GitHubRepos: 10,
		AWSAccounts: 5,
		TeamMembers: 10,
		CustomRules: 50,
		APIRequests: 10000,
	},
	"enterprise": {
		GitHubRepos: -1, // Unlimited
		AWSAccounts: -1, // Unlimited
		TeamMembers: -1, // Unlimited
		CustomRules: -1, // Unlimited
		APIRequests: -1, // Unlimited
	},
}

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Plan           string `json:"plan" binding:"required,oneof=free pro enterprise"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	Status  string    `json:"status" binding:"required,oneof=active inactive suspended cancelled"`
	EndDate time.Time `json:"end_date" binding:"required"`
}
