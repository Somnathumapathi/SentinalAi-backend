package main

import (
	"log"
	"net/http"
	"os"
	"sentinal/controller"
	"sentinal/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "sentinal/controller/razorpay"
)

// var db = make(map[string]string)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Debug: Print environment variables
	log.Printf("SUPABASE_URL: %s", os.Getenv("SUPABASE_URL"))
	log.Printf("SUPABASE_KEY: %s", os.Getenv("SUPABASE_KEY"))
	log.Printf("PORT: %s", os.Getenv("PORT"))

	// Initialize database connection with GORM
	if err := db.InitGORM(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set up Gin router
	r := gin.Default()

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
