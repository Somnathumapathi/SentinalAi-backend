package models

import (
	"time"
)

// GitHubInstallation represents a GitHub App installation
type GitHubInstallation struct {
	ID                  string            `json:"id" bson:"_id"`
	UserID              string            `json:"user_id" bson:"user_id"`
	OrganizationID      string            `json:"organization_id" bson:"organization_id"`
	InstallationID      int64             `json:"installation_id" bson:"installation_id"`
	AccountID           int64             `json:"account_id" bson:"account_id"`
	AccountType         string            `json:"account_type" bson:"account_type"`
	AccountLogin        string            `json:"account_login" bson:"account_login"`
	RepositorySelection string            `json:"repository_selection" bson:"repository_selection"`
	AccessTokensURL     string            `json:"access_tokens_url" bson:"access_tokens_url"`
	RepositoriesURL     string            `json:"repositories_url" bson:"repositories_url"`
	HTMLURL             string            `json:"html_url" bson:"html_url"`
	AppID               int64             `json:"app_id" bson:"app_id"`
	TargetID            int64             `json:"target_id" bson:"target_id"`
	TargetType          string            `json:"target_type" bson:"target_type"`
	Permissions         map[string]string `json:"permissions" bson:"permissions"`
	Events              []string          `json:"events" bson:"events"`
	CreatedAt           time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at" bson:"updated_at"`
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	ID             string    `json:"id" bson:"_id"`
	InstallationID string    `json:"installation_id" bson:"installation_id"`
	RepoID         int64     `json:"repo_id" bson:"repo_id"`
	Name           string    `json:"name" bson:"name"`
	FullName       string    `json:"full_name" bson:"full_name"`
	Private        bool      `json:"private" bson:"private"`
	HTMLURL        string    `json:"html_url" bson:"html_url"`
	Description    string    `json:"description" bson:"description"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
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
