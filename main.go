package main

import (
	"log"
	"net/http"
	"os"
	"sentinal/controller"
	"sentinal/db"
	"sentinal/middleware"
	"sentinal/webhooks"

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

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/logout", controller.SignOut)
		auth.GET("/user", controller.GetUser)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// GitHub integration routes
		protected.GET("/github/app/install-url", controller.GetGitHubAppInstallURL)
		protected.GET("/github/installations", controller.GetGitHubInstallations)
		protected.POST("/github/installations/:id/connect", controller.ConnectGitHubInstallation)

		// Subscription routes
		protected.GET("/subscriptions/:org_id", controller.GetSubscription)
		protected.POST("/subscriptions", controller.CreateSubscription)
		protected.PUT("/subscriptions/:id", controller.UpdateSubscription)
		protected.DELETE("/subscriptions/:id", controller.CancelSubscription)
	}

	// GitHub webhook endpoint
	r.POST("/webhook/github", webhooks.GitHubWebhook)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
