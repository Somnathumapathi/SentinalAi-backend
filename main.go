package main

import (
	"net/http"
	"sentinal/controller"

	"github.com/gin-gonic/gin"
	// "sentinal/controller/razorpay"
)

// var db = make(map[string]string)

func main() {
	// Load environment variables

	// Set up Gin router
	r := gin.Default()
	// handle cors
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/trace", controller.TraceHandler)
	r.POST("/webhook/github", controller.GitHubIWebhook)
	r.POST("/getiac", controller.GetIacContent)
	r.POST("/createpr", controller.CreatePRHandler)

	r.Run(":8080")
}

//
