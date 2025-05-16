package webhooks

// import (
// 	"fmt"
// 	"net/http"
// "github.com/gin-gonic/gin"
// )

// func GitHubIWebhook(c *gin.Context) {
// 	// Parse the request body
// 	var traceRequest TraceRequest
// 	if err := c.ShouldBindJSON(&traceRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Process the misconfiguration
// 	err := processMisConfig(traceRequest)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
// }
