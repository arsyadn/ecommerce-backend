package controllers

import (
	"final-project/models"
	"final-project/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	OrderService *services.OrderService
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		OrderService: services.NewOrderService(db),
	}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order []models.OrderPayload
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.GetUint("user_id")
	role, err := oc.OrderService.GetUserRole(userID)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user role"})
		return
	}

	if role != "user" {
		c.JSON(403, gin.H{"error": "Only user can make order"})
		return
	}


	for _, detail := range order {
		if detail.Quantity <= 0 {
			c.JSON(400, gin.H{"error": "Quantity must be greater than 0"})
			return
		}
		if detail.ItemID == 0 {
			c.JSON(400, gin.H{"error": "Item ID is required"})
			return
		}
		if detail.PriceAtOrder <= 0 {
			c.JSON(400, gin.H{"error": "Price must be greater than 0"})
			return
		}
	}

	user := c.GetUint("user_id")
	newOrder := models.Order{
		UserID: user,
		Status: "pending",
	}

	// Convert payload 
	var orderDetails []models.OrderDetail
	for _, item := range order {
		detail := models.OrderDetail{
			ItemID:       item.ItemID,
			Quantity:     item.Quantity,
			PriceAtOrder: float64(item.PriceAtOrder),
		}
		orderDetails = append(orderDetails, detail)
	}
	newOrder.Details = orderDetails

	if err := oc.OrderService.CreateOrder(&newOrder); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Order created successfully", "status": "success"})
}