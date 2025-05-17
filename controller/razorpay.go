package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func GetRazorpayOrder(c *gin.Context) {
	key_id := os.Getenv("KEY_ID")
	key_secret := os.Getenv("KEY_SECRET")
	client := razorpay.NewClient(key_id, key_secret)

	orderData := map[string]interface{}{
		"amount":   50000,
		"currency": "INR",
		"receipt":  "receipt#1",
	}
	order, err := client.Order.Create(orderData, nil)
	if err != nil {
		fmt.Println("Error creating order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Order ID:", order["id"])
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully"})
}
