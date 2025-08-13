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
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if len(order.Details) == 0 {
		c.JSON(400, gin.H{"error": "Order must include details"})
		return
	}

	for _, detail := range order.Details {
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

	order.UserID = c.GetUint("user_id")

	if err := oc.OrderService.CreateOrder(&order); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Order created successfully", "status": "success"})
}