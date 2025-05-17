package models

import "time"

type TraceRequest struct {
	resource     string `json:"resource"`
	misconfig    string `json:"misconfig"`
	account      string `json:"account"`
	organization string `json:"organization"`
}

type GitHubIWebhook struct {
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` // Hashed password
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// type Organization struct {
// 	ID        string    `json:"id" bson:"_id,omitempty"`
// 	Name      string    `json:"name" bson:"name"`
// 	OwnerID   string    `json:"owner_id" bson:"owner_id"`
// 	CreatedAt time.Time `json:"created_at" bson:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
// }

// type Subscription struct {
// 	ID                     string    `json:"id" bson:"_id,omitempty"`
// 	OrganizationID         string    `json:"organization_id" bson:"organization_id"`
// 	Plan                   string    `json:"plan" bson:"plan"`     // free, pro, enterprise
// 	Status                 string    `json:"status" bson:"status"` // active, inactive, suspended
// 	RazorpaySubscriptionID string    `json:"razorpay_subscription_id" bson:"razorpay_subscription_id"`
// 	StartDate              time.Time `json:"start_date" bson:"start_date"`
// 	EndDate                time.Time `json:"end_date" bson:"end_date"`
// 	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
// 	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
// }

type GitHubIntegration struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	OrganizationID string    `json:"organization_id" bson:"organization_id"`
	InstallationID int64     `json:"installation_id" bson:"installation_id"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

type AWSIntegration struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	OrganizationID  string    `json:"organization_id" bson:"organization_id"`
	AccessKeyID     string    `json:"access_key_id" bson:"access_key_id"`
	SecretAccessKey string    `json:"secret_access_key" bson:"secret_access_key"`
	Region          string    `json:"region" bson:"region"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// PlanLimits defines the limits for each subscription plan
// type PlanLimits struct {
// 	GitHubRepos int `json:"github_repos" bson:"github_repos"`
// 	AWSAccounts int `json:"aws_accounts" bson:"aws_accounts"`
// 	TeamMembers int `json:"team_members" bson:"team_members"`
// 	CustomRules int `json:"custom_rules" bson:"custom_rules"`
// 	APIRequests int `json:"api_requests" bson:"api_requests"`
// }

// // PlanLimitsMap defines the limits for each plan type
// var PlanLimitsMap = map[string]PlanLimits{
// 	"free": {
// 		GitHubRepos: 1,
// 		AWSAccounts: 1,
// 		TeamMembers: 2,
// 		CustomRules: 5,
// 		APIRequests: 1000,
// 	},
// 	"pro": {
// 		GitHubRepos: 10,
// 		AWSAccounts: 5,
// 		TeamMembers: 10,
// 		CustomRules: 50,
// 		APIRequests: 10000,
// 	},
// 	"enterprise": {
// 		GitHubRepos: -1, // Unlimited
// 		AWSAccounts: -1, // Unlimited
// 		TeamMembers: -1, // Unlimited
// 		CustomRules: -1, // Unlimited
// 		APIRequests: -1, // Unlimited
// 	},
// }

// Request/Response models
// type CreateSubscriptionRequest struct {
// 	OrganizationID string `json:"organization_id" binding:"required"`
// 	Plan           string `json:"plan" binding:"required"`
// }

// type UpdateSubscriptionRequest struct {
// 	Status  string    `json:"status"`
// 	EndDate time.Time `json:"end_date"`
// }

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateGitHubIntegrationRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	InstallationID int64  `json:"installation_id" binding:"required"`
}

type CreateAWSIntegrationRequest struct {
	OrganizationID  string `json:"organization_id" binding:"required"`
	AccessKeyID     string `json:"access_key_id" binding:"required"`
	SecretAccessKey string `json:"secret_access_key" binding:"required"`
	Region          string `json:"region" binding:"required"`
}
