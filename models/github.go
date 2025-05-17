package models

import (
	"time"

	"gorm.io/gorm"
)

// GitHubInstallation represents a GitHub App installation
type GitHubInstallation struct {
	ID                  string            `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID              string            `json:"user_id" gorm:"not null"`
	OrganizationID      string            `json:"organization_id" gorm:"not null"`
	InstallationID      int64             `json:"installation_id" gorm:"not null;uniqueIndex"`
	AccountID           int64             `json:"account_id" gorm:"not null"`
	AccountType         string            `json:"account_type" gorm:"not null"`
	AccountLogin        string            `json:"account_login" gorm:"not null"`
	RepositorySelection string            `json:"repository_selection" gorm:"not null"`
	AccessTokensURL     string            `json:"access_tokens_url" gorm:"not null"`
	RepositoriesURL     string            `json:"repositories_url" gorm:"not null"`
	HTMLURL             string            `json:"html_url" gorm:"not null"`
	AppID               int64             `json:"app_id" gorm:"not null"`
	TargetID            int64             `json:"target_id" gorm:"not null"`
	TargetType          string            `json:"target_type" gorm:"not null"`
	Permissions         map[string]string `json:"permissions" gorm:"type:jsonb;not null"`
	Events              []string          `json:"events" gorm:"type:jsonb;not null"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	DeletedAt           gorm.DeletedAt    `json:"-" gorm:"index"`
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	InstallationID string         `json:"installation_id" gorm:"not null"`
	RepoID         int64          `json:"repo_id" gorm:"not null"`
	Name           string         `json:"name" gorm:"not null"`
	FullName       string         `json:"full_name" gorm:"not null"`
	Private        bool           `json:"private" gorm:"not null"`
	HTMLURL        string         `json:"html_url" gorm:"not null"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// GitHubWebhookPayload represents the payload received from GitHub webhooks
type GitHubWebhookPayload struct {
	Action       string `json:"action"`
	Installation struct {
		ID                  int64             `json:"id"`
		AccountID           int64             `json:"account_id"`
		AccountType         string            `json:"account_type"`
		AccountLogin        string            `json:"account_login"`
		RepositorySelection string            `json:"repository_selection"`
		AccessTokensURL     string            `json:"access_tokens_url"`
		RepositoriesURL     string            `json:"repositories_url"`
		HTMLURL             string            `json:"html_url"`
		AppID               int64             `json:"app_id"`
		TargetID            int64             `json:"target_id"`
		TargetType          string            `json:"target_type"`
		Permissions         map[string]string `json:"permissions"`
		Events              []string          `json:"events"`
	} `json:"installation"`
	Repositories []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Private     bool   `json:"private"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
	} `json:"repositories"`
	Repository struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Private     bool   `json:"private"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
	} `json:"repository"`
}

// CreateGitHubIntegrationRequest represents the request to create a GitHub integration
type CreateGitHubIntegrationRequest struct {
	OrganizationID string `json:"organization_id" binding:"required"`
	InstallationID int64  `json:"installation_id" binding:"required"`
}

// TraceRequest represents a request to trace infrastructure configuration
type TraceRequest struct {
	RepositoryFullName string `json:"repository_full_name" binding:"required"`
	PullRequestNumber  int    `json:"pull_request_number" binding:"required"`
}

// GitHubIWebhook represents a GitHub installation webhook payload
type GitHubIWebhook struct {
	Action       string `json:"action"`
	Installation struct {
		ID                  int64             `json:"id"`
		AccountID           int64             `json:"account_id"`
		AccountType         string            `json:"account_type"`
		AccountLogin        string            `json:"account_login"`
		RepositorySelection string            `json:"repository_selection"`
		AccessTokensURL     string            `json:"access_tokens_url"`
		RepositoriesURL     string            `json:"repositories_url"`
		HTMLURL             string            `json:"html_url"`
		AppID               int64             `json:"app_id"`
		TargetID            int64             `json:"target_id"`
		TargetType          string            `json:"target_type"`
		Permissions         map[string]string `json:"permissions"`
		Events              []string          `json:"events"`
	} `json:"installation"`
	Repository struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Private     bool   `json:"private"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
	} `json:"repository"`
}
