package controller

import (
	"net/http"
	"sentinal/db"
	"sentinal/models"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Register user with Supabase Auth
	user, err := db.SignUp(c.Request.Context(), req.Email, req.Password, map[string]interface{}{
		"name": req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Create user profile using GORM
	profile := models.User{
		ID:        user.ID.String(),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user profile to database using GORM
	if err := db.GetDB().Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user":    user,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login with Supabase Auth
	session, err := db.SignIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"session": session,
	})
}

// GetUser retrieves the current user's information
func GetUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		return
	}

	// Get user from Supabase Auth
	user, err := db.GetUser(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get user profile using GORM
	var profile models.User
	if err := db.GetDB().Where("id = ?", user.ID.String()).First(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"profile": profile,
	})
}

// SignOut handles user logout
func SignOut(c *gin.Context) {
	// Note: Supabase handles token invalidation automatically
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}
