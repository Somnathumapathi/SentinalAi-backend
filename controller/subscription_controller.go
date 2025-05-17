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
)

// GetSubscription retrieves the subscription for an organization
func GetSubscription(c *gin.Context) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var subscription models.Subscription
	err := db.Select("subscriptions", &subscription, "organization_id", orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription"})
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

	// Create subscription
	subscription := models.Subscription{
		OrganizationID:         req.OrganizationID,
		Plan:                   req.Plan,
		Status:                 "pending",
		RazorpaySubscriptionID: "", // Will be set after Razorpay subscription is created
		StartDate:              time.Now(),
		EndDate:                time.Now().AddDate(1, 0, 0), // 1 year subscription
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	err := db.Insert("subscriptions", subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
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

	update := map[string]interface{}{
		"status":     req.Status,
		"end_date":   req.EndDate,
		"updated_at": time.Now(),
	}

	err := db.Update("subscriptions", update, "id", subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
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

	err := db.Update("subscriptions", map[string]interface{}{
		"status":     "cancelled",
		"updated_at": time.Now(),
	}, "id", subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

func createRazorpaySubscription(plan string) (map[string]interface{}, error) {
	key_id := os.Getenv("KEY_ID")
	key_secret := os.Getenv("KEY_SECRET")
	client := razorpay.NewClient(key_id, key_secret)

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
	err := db.Select("subscriptions", &subscription, "organization_id", orgID, "status", "active")
	if err != nil {
		return false, err
	}

	limits := models.PlanLimitsMap[subscription.Plan]

	// Check GitHub repos limit
	if limits.GitHubRepos != -1 {
		var count int
		err := db.Select("github_integrations", &count, "organization_id", orgID)
		if err != nil {
			return false, err
		}
		if count >= limits.GitHubRepos {
			return true, nil
		}
	}

	// Check AWS accounts limit
	if limits.AWSAccounts != -1 {
		var count int
		err := db.Select("aws_integrations", &count, "organization_id", orgID)
		if err != nil {
			return false, err
		}
		if count >= limits.AWSAccounts {
			return true, nil
		}
	}

	return false, nil
}
