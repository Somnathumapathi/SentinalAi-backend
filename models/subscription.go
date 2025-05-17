package models

import (
	"time"
)

// PlanLimits defines the limits for each subscription plan
type PlanLimits struct {
	GitHubRepos int `json:"github_repos" bson:"github_repos"`
	AWSAccounts int `json:"aws_accounts" bson:"aws_accounts"`
	TeamMembers int `json:"team_members" bson:"team_members"`
	CustomRules int `json:"custom_rules" bson:"custom_rules"`
	APIRequests int `json:"api_requests" bson:"api_requests"`
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

// Subscription represents a subscription in the system
type Subscription struct {
	ID                     string    `json:"id" bson:"_id"`
	OrganizationID         string    `json:"organization_id" bson:"organization_id"`
	Plan                   string    `json:"plan" bson:"plan"`
	Status                 string    `json:"status" bson:"status"`
	RazorpaySubscriptionID string    `json:"razorpay_subscription_id" bson:"razorpay_subscription_id"`
	StartDate              time.Time `json:"start_date" bson:"start_date"`
	EndDate                time.Time `json:"end_date" bson:"end_date"`
	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateSubscriptionRequest represents the request to create a subscription
type CreateSubscriptionRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	Plan           string `json:"plan" binding:"required,oneof=free pro enterprise"`
}

// UpdateSubscriptionRequest represents the request to update a subscription
type UpdateSubscriptionRequest struct {
	Status  string    `json:"status" binding:"required,oneof=active cancelled expired"`
	EndDate time.Time `json:"end_date" binding:"required"`
}

// Organization represents an organization in the system
type Organization struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	OwnerID   string    `json:"owner_id" bson:"owner_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
