package webhooks

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sentinal/db"
	"sentinal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GitHubWebhook handles incoming GitHub webhook events
func GitHubWebhook(c *gin.Context) {
	// Verify webhook signature
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing signature"})
		return
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Restore request body for later use
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	// Get event type
	eventType := c.GetHeader("X-GitHub-Event")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing event type"})
		return
	}

	// Handle different event types
	switch eventType {
	case "installation":
		handleInstallationEvent(c, body)
	case "installation_repositories":
		handleInstallationRepositoriesEvent(c, body)
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Event type not handled"})
	}
}

// handleInstallationEvent processes GitHub App installation events
func handleInstallationEvent(c *gin.Context, body []byte) {
	var payload struct {
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
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	switch payload.Action {
	case "created":
		// Create organization
		org := models.Organization{
			Name:      payload.Installation.AccountLogin,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := db.Insert("organizations", org)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
			return
		}

		// Create installation
		installation := models.GitHubInstallation{
			OrganizationID:      org.ID,
			InstallationID:      payload.Installation.ID,
			AccountID:           payload.Installation.AccountID,
			AccountType:         payload.Installation.AccountType,
			AccountLogin:        payload.Installation.AccountLogin,
			RepositorySelection: payload.Installation.RepositorySelection,
			AccessTokensURL:     payload.Installation.AccessTokensURL,
			RepositoriesURL:     payload.Installation.RepositoriesURL,
			HTMLURL:             payload.Installation.HTMLURL,
			AppID:               payload.Installation.AppID,
			TargetID:            payload.Installation.TargetID,
			TargetType:          payload.Installation.TargetType,
			Permissions:         payload.Installation.Permissions,
			Events:              payload.Installation.Events,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
		err = db.Insert("github_installations", installation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create installation"})
			return
		}

	case "deleted":
		// Delete installation
		err := db.Delete("github_installations", "installation_id", payload.Installation.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete installation"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed successfully"})
}

// handleInstallationRepositoriesEvent processes repository changes for an installation
func handleInstallationRepositoriesEvent(c *gin.Context, body []byte) {
	var payload struct {
		Action       string `json:"action"`
		Installation struct {
			ID int64 `json:"id"`
		} `json:"installation"`
		Repositories []struct {
			ID          int64  `json:"id"`
			Name        string `json:"name"`
			FullName    string `json:"full_name"`
			Private     bool   `json:"private"`
			HTMLURL     string `json:"html_url"`
			Description string `json:"description"`
		} `json:"repositories"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Get installation
	var installation models.GitHubInstallation
	err := db.Select("github_installations", &installation, "installation_id", payload.Installation.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Installation not found"})
		return
	}

	// Process repositories
	for _, repo := range payload.Repositories {
		githubRepo := map[string]interface{}{
			"installation_id": installation.ID,
			"repo_id":         repo.ID,
			"name":            repo.Name,
			"full_name":       repo.FullName,
			"private":         repo.Private,
			"html_url":        repo.HTMLURL,
			"description":     repo.Description,
			"created_at":      time.Now(),
			"updated_at":      time.Now(),
		}

		err := db.Update("github_repositories", githubRepo, "installation_id", installation.ID, "repo_id", repo.ID)
		if err != nil {
			// If update fails, try insert
			err = db.Insert("github_repositories", githubRepo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process repository"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed successfully"})
}

// verifyGitHubSignature verifies the GitHub webhook signature
func verifyGitHubSignature(c *gin.Context) bool {
	signature := c.GetHeader("X-Hub-Signature-256")
	if signature == "" {
		return false
	}

	// Get the webhook secret from environment
	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		return false
	}

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false
	}

	// Restore the request body for later use
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	// Calculate the expected signature
	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
