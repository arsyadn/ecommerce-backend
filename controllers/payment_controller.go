package controllers

import (
	"final-project/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentController struct {
	PaymentService *services.PaymentService
}

func NewPaymentController(db *gorm.DB) *PaymentController {
	return &PaymentController{
		PaymentService: services.NewPaymentService(db),
	}
}


func (pc *PaymentController) UploadPayment(c *gin.Context) {
    userID := c.GetUint("user_id")
    role, err := pc.PaymentService.GetUserRole(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to get user role"})
        return
    }

    if role != "user" {
        c.JSON(403, gin.H{"error": "Only user can make payment"})
        return
    }

    // Read form-data values
    orderID, err := strconv.Atoi(c.PostForm("order_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order_id"})
        return
    }

    amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
        return
    }

    // Read file
    file, _, err := c.Request.FormFile("evidence")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Evidence file is required"})
        return
    }
    defer file.Close()

    // Call service
    if err := pc.PaymentService.UploadPayment(
        c.Request.Context(),
        uint(orderID),
        amount,
        file,
    ); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Payment uploaded successfully"})
}

// GetPayments returns all payments
func (pc *PaymentController) GetPayments(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, err := pc.PaymentService.GetUserRole(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user role"})
		return
	}

	// Check if user is admin
	if role != "admin" {
		c.JSON(403, gin.H{"error": "Only admin can create items"})
		return
	}

	payments, err := pc.PaymentService.GetPayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}
