package controller

import (
	"fmt"
	"net/http"
	"os"
	"sentinal/db"
	"sentinal/models"
	"time"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

// GetSubscription retrieves the subscription for an organization
func GetSubscription(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var subscription models.Subscription
	result := db.DB.Where("organization_id = ?", orgID).First(&subscription)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch subscription: %v", result.Error)})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// CreateSubscription creates a new subscription for an organization
func CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate plan
	if _, exists := models.PlanLimitsMap[req.Plan]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan"})
		return
	}

	// Check if organization already has an active subscription
	var existingSubscription models.Subscription
	result := db.DB.Where("organization_id = ? AND status = ?", req.OrganizationID, "active").First(&existingSubscription)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization already has an active subscription"})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check existing subscription: %v", result.Error)})
		return
	}

	// Create Razorpay subscription
	razorpaySub, err := createRazorpaySubscription(req.Plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Razorpay subscription: %v", err)})
		return
	}

	// Create subscription
	subscription := models.Subscription{
		OrganizationID:         req.OrganizationID,
		Plan:                   req.Plan,
		Status:                 "pending",
		RazorpaySubscriptionID: razorpaySub["id"].(string),
		StartDate:              time.Now(),
		EndDate:                time.Now().AddDate(1, 0, 0), // 1 year subscription
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	result = db.DB.Create(&subscription)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create subscription: %v", result.Error)})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// UpdateSubscription updates an existing subscription
func UpdateSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":    true,
		"inactive":  true,
		"suspended": true,
		"cancelled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	result := db.DB.Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"status":     req.Status,
			"end_date":   req.EndDate,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update subscription: %v", result.Error)})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// CancelSubscription cancels a subscription
func CancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	result := db.DB.Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"status":     "cancelled",
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to cancel subscription: %v", result.Error)})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

func createRazorpaySubscription(plan string) (map[string]interface{}, error) {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	if keyID == "" || keySecret == "" {
		return nil, fmt.Errorf("Razorpay credentials not configured")
	}

	client := razorpay.NewClient(keyID, keySecret)

	// Define plan prices
	planPrices := map[string]int64{
		"free":       0,      // Free tier
		"pro":        150000, // ₹1,500/month
		"enterprise": 500000, // ₹5,000/month
	}

	amount, exists := planPrices[plan]
	if !exists {
		return nil, fmt.Errorf("invalid plan: %s", plan)
	}

	// Create a plan first
	planData := map[string]interface{}{
		"period":   "monthly",
		"interval": 1,
		"item": map[string]interface{}{
			"name":     fmt.Sprintf("SentinelAI %s Plan", plan),
			"amount":   amount,
			"currency": "INR",
		},
	}

	createdPlan, err := client.Plan.Create(planData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %v", err)
	}

	// Create subscription
	subscriptionData := map[string]interface{}{
		"plan_id":         createdPlan["id"],
		"customer_notify": 1,
		"total_count":     12, // 12 months
	}

	subscription, err := client.Subscription.Create(subscriptionData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %v", err)
	}

	return subscription, nil
}

// Helper function to check if an organization has reached its plan limits
func CheckPlanLimits(orgID string) (bool, error) {
	var subscription models.Subscription
	result := db.DB.Where("organization_id = ? AND status = ?", orgID, "active").First(&subscription)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("no active subscription found")
		}
		return false, fmt.Errorf("failed to fetch subscription: %v", result.Error)
	}

	limits := models.PlanLimitsMap[subscription.Plan]

	// Check GitHub repos limit
	if limits.GitHubRepos != -1 {
		var count int64
		result := db.DB.Model(&models.GitHubInstallation{}).Where("organization_id = ?", orgID).Count(&count)
		if result.Error != nil {
			return false, fmt.Errorf("failed to count GitHub installations: %v", result.Error)
		}
		if count >= int64(limits.GitHubRepos) {
			return true, nil
		}
	}

	// Check AWS accounts limit
	if limits.AWSAccounts != -1 {
		var count int64
		result := db.DB.Model(&models.AWSIntegration{}).Where("organization_id = ?", orgID).Count(&count)
		if result.Error != nil {
			return false, fmt.Errorf("failed to count AWS integrations: %v", result.Error)
		}
		if count >= int64(limits.AWSAccounts) {
			return true, nil
		}
	}

	return false, nil
}
