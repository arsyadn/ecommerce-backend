package controllers

import (
	"final-project/models"
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
		c.JSON(403, gin.H{"error": "Only admin can see payments"})
		return
	}

	payments, err := pc.PaymentService.GetPayments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (pc *PaymentController) GetPaymentByID(c *gin.Context) {
    paymentID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
        return
    }

    payment, err := pc.PaymentService.GetPaymentByID(c.Request.Context(), paymentID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, payment)
}

func (pc *PaymentController) AdminUpdatePayment(c *gin.Context) {
    userID := c.GetUint("user_id")
    role, err := pc.PaymentService.GetUserRole(userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to get user role"})
        return
    }

    if role != "admin" {
        c.JSON(403, gin.H{"error": "Only admin can update payment"})
        return
    }
    // Read JSON payload
    var payload models.UpdatePaymentStatusPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    paymentID := payload.PaymentID
    status := payload.Status

    // Call service
    if err := pc.PaymentService.AdminUpdatePayment(c.Request.Context(), paymentID, status, int(userID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully"})
}