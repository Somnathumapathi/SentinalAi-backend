package controller

import (
	"fmt"
	"net/http"
	"os"
	"sentinal/models"
	githubsvc "sentinal/services/github"
	"strings"
	"time"

	"sentinal/db"

	"github.com/gin-gonic/gin"
	github "github.com/google/go-github/v53/github"
)

// TraceHandler handles infrastructure configuration tracing requests
func TraceHandler(c *gin.Context) {
	var traceRequest models.TraceRequest
	if err := c.ShouldBindJSON(&traceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go processMisConfig(c, traceRequest)
}

// GitHubIWebhook handles GitHub installation webhook events
func GitHubIWebhook(c *gin.Context) {
	var webhook models.GitHubIWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error binding JSON: %v", err)})
		return
	}

	// Process the webhook based on the action
	switch webhook.Action {
	case "created":
		// Handle new installation
		if err := handleNewInstallation(c, webhook); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to handle new installation: %v", err)})
			return
		}
	case "deleted":
		// Handle installation deletion
		if err := handleInstallationDeletion(c, webhook); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to handle installation deletion: %v", err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func processMisConfig(c *gin.Context, req models.TraceRequest) {
	// Get GitHub client for the installation
	// TODO: Get installation ID and app ID from the database based on repository
	installationID := int64(0) // Replace with actual installation ID
	appID := int64(0)          // Replace with actual app ID
	client, err := githubsvc.GetGHClient(installationID, appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get GitHub client: %v", err)})
		return
	}

	// Get repository owner and name from full name
	parts := strings.Split(req.RepositoryFullName, "/")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository full name format"})
		return
	}
	owner, repo := parts[0], parts[1]

	// Get changed files in the PR
	files, _, err := client.PullRequests.ListFiles(c, owner, repo, req.PullRequestNumber, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to list PR files: %v", err)})
		return
	}

	// Process each changed file
	for _, file := range files {
		if strings.HasSuffix(file.GetFilename(), ".tf") {
			content, err := getDecodedFileContent(c, client, owner, repo, file.GetFilename())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get file content: %v", err)})
				return
			}
			// Process the Terraform file content
			// TODO: Add Terraform file processing logic
			fmt.Printf("Processing Terraform file: %s\nContent length: %d\n", file.GetFilename(), len(content))
		}
	}
}

// GetGitHubAppInstallURL returns the URL for installing the GitHub App
func GetGitHubAppInstallURL(c *gin.Context) {
	appID := os.Getenv("GITHUB_APP_ID")
	if appID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub App ID not configured"})
		return
	}

	installURL := fmt.Sprintf("https://github.com/apps/%s/installations/new", appID)
	c.JSON(http.StatusOK, gin.H{"install_url": installURL})
}

// GetGitHubInstallations retrieves all installations for a user's organizations
func GetGitHubInstallations(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var installations []models.GitHubInstallation
	result := db.DB.Where("user_id = ?", userID).Find(&installations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch installations: %v", result.Error)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"installations": installations})
}

// ConnectGitHubInstallation connects a GitHub installation to a user's organization
func ConnectGitHubInstallation(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateGitHubIntegrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the installation with the user ID
	result := db.DB.Model(&models.GitHubInstallation{}).
		Where("installation_id = ?", req.InstallationID).
		Updates(map[string]interface{}{
			"user_id":    userID,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update installation: %v", result.Error)})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Installation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation connected successfully"})
}

// Helper functions
func handleNewInstallation(c *gin.Context, webhook models.GitHubIWebhook) error {
	installation := models.GitHubInstallation{
		InstallationID:      webhook.Installation.ID,
		AccountID:           webhook.Installation.AccountID,
		AccountType:         webhook.Installation.AccountType,
		AccountLogin:        webhook.Installation.AccountLogin,
		RepositorySelection: webhook.Installation.RepositorySelection,
		AccessTokensURL:     webhook.Installation.AccessTokensURL,
		RepositoriesURL:     webhook.Installation.RepositoriesURL,
		HTMLURL:             webhook.Installation.HTMLURL,
		AppID:               webhook.Installation.AppID,
		TargetID:            webhook.Installation.TargetID,
		TargetType:          webhook.Installation.TargetType,
		Permissions:         webhook.Installation.Permissions,
		Events:              webhook.Installation.Events,
	}

	result := db.DB.Create(&installation)
	return result.Error
}

func handleInstallationDeletion(c *gin.Context, webhook models.GitHubIWebhook) error {
	result := db.DB.Where("installation_id = ?", webhook.Installation.ID).Delete(&models.GitHubInstallation{})
	return result.Error
}

func getDecodedFileContent(ctx *gin.Context, client *github.Client, owner, repo, filePath string) (string, error) {
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get file content: %v", err)
	}

	if fileContent == nil {
		return "", fmt.Errorf("file content is nil for path: %s", filePath)
	}

	decoded, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %v", err)
	}

	return decoded, nil
}
