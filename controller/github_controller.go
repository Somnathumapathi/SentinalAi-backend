package controller

import (
	"fmt"
	// "fmt"
	"net/http"
	"os"
	"sentinal/models"
	githubsvc "sentinal/services/github"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	github "github.com/google/go-github/v53/github"

	// "github.com/joho/godotenv"
	"sentinal/db"
)

func TraceHandler(c *gin.Context) {
	var traceRequest models.TraceRequest
	if err := c.ShouldBindJSON(&traceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	go processMisConfig(c, traceRequest)

}

func GitHubIWebhook(c *gin.Context) {
	// Parse the request body
	var githubIWebhook models.GitHubIWebhook
	if err := c.BindJSON(&githubIWebhook); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// installationId := githubIWebhook.Installation.ID
	// repoFullName := githubIWebhook.Repository.FullName
	// You can now use the installationId and repoFullName to perform actions
	getIaCFileContent(c)
	// fmt.Println("Installation ID:", installationId)
	// fmt.Println("Repository Full Name:", repoFullName)

}

func processMisConfig(c *gin.Context, req models.TraceRequest) {
	fmt.Println("Reached")
	client, _ := githubsvc.GetGHClient(0000000, 0000000)
	fmt.Println("Client:", client)
	//find the pr
	prs, _, err := client.PullRequests.ListFiles(c, "Somnathumapathi", "CraveHub", 10, nil)
	if err != nil {
		fmt.Println("Error listing pull requests:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, pr := range prs {
		fmt.Println("PR:", pr)
	}

}

func getIaCFileContent(c *gin.Context) {
	client, err := githubsvc.GetGHClient(int64(67171708), int64(1271564))
	if err != nil {
		fmt.Printf("Error getting GitHub client: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize GitHub client"})
		return
	}
	if client == nil {
		fmt.Println("GitHub client is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub client is nil"})
		return
	}
	collectIaCFiles(c, client, "Somnathumapathi", "CraveHub", "", []string{".tf"})
}

func collectIaCFiles(ctx *gin.Context, client *github.Client, owner, repo, path string, extensions []string) {
	if client == nil {
		fmt.Println("GitHub client is nil in collectIaCFiles")
		return
	}

	fileContent, dirContents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		fmt.Printf("Error getting contents at path %s: %v\n", path, err)
		return
	}

	// If dirContents is not nil, it's a directory
	if dirContents != nil {
		for _, content := range dirContents {
			if content == nil {
				continue
			}
			fmt.Print("Content: ", content)
			switch content.GetType() {
			case "file":
				for _, ext := range extensions {
					if strings.HasSuffix(content.GetPath(), ext) {
						fmt.Printf("Found IaC file: %s\n", content.GetPath())
						fileContent, err := getDecodedFileContent(ctx, client, owner, repo, content.GetPath())
						if err != nil {
							fmt.Printf("Error decoding %s: %v\n", content.GetPath(), err)
							continue
						}
						fmt.Println("File content:", fileContent[:min(300, len(fileContent))])
					}
				}
			case "dir":
				collectIaCFiles(ctx, client, owner, repo, content.GetPath(), extensions)
			}
		}
		return
	}

	// If fileContent is not nil, it's a file
	if fileContent != nil {
		for _, ext := range extensions {
			if strings.HasSuffix(fileContent.GetPath(), ext) {
				fmt.Printf("Found IaC file: %s\n", fileContent.GetPath())
				decoded, err := fileContent.GetContent()
				if err != nil {
					fmt.Printf("Error decoding %s: %v\n", fileContent.GetPath(), err)
					return
				}
				fmt.Println("File content:", decoded[:min(300, len(decoded))])
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getDecodedFileContent(ctx *gin.Context, client *github.Client, owner, repo, filePath string) (string, error) {
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
	if err != nil {
		return "", err
	}

	if fileContent == nil {
		return "", fmt.Errorf("file content is nil for path: %s", filePath)
	}

	decoded, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}

	return decoded, nil
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
	err := db.Select("github_installations", &installations, "user_id", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch installations"})
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

	installationID := c.Param("id")
	if installationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Installation ID is required"})
		return
	}

	// Get the installation from the database
	var installation models.GitHubInstallation
	err := db.Select("github_installations", &installation, "installation_id", installationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Installation not found"})
		return
	}

	// Update the installation with the user ID
	update := map[string]interface{}{
		"user_id":    userID,
		"updated_at": time.Now(),
	}

	err = db.Update("github_installations", update, "installation_id", installationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update installation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation connected successfully"})
}
