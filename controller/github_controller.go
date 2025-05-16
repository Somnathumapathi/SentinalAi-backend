package controller

import (
	"fmt"
	"log"
	"net/http"
	"sentinal/models"
	"sentinal/services/github"

	"github.com/gin-gonic/gin"
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
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	installationId := githubIWebhook.Installation.ID
	repoFullName := githubIWebhook.Repository.FullName
	log.Println("Installation ID:", installationId)
	log.Println("Repository Full Name:", repoFullName)
}

func processMisConfig(c *gin.Context, req models.TraceRequest) {
	fmt.Println("Reached")
	client, _ := github.GetGHClient(0000000, 0000000)
	fmt.Println("Client:", client)
	//find the pr
	prs, _, err := client.PullRequests.ListFiles(c, "Somnathumapathi", "CraveHub", 10, nil)
	if err != nil {
		log.Println("Error listing pull requests:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, pr := range prs {
		fmt.Println("PR:", pr)
	}

}
