package controllers

import (
	"final-project/models"
	"final-project/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemController struct {
	ItemService *services.ItemService
}

func NewItemController(db *gorm.DB) *ItemController {
	return &ItemController{
		ItemService: services.NewItemService(db),
	}
}


func (ic *ItemController) CreateItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	item.UserID = userID

	if err := ic.ItemService.CreateItem(&item); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(201, gin.H{"message": "Item created successfully", "status": "success"})
}


func (ic *ItemController) GetAllItems(c *gin.Context) {
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	items, err := ic.ItemService.GetAllItems(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve items"})
		return
	}

	c.JSON(200, gin.H{
		"items": items,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
		},
	})
}

func (ic *ItemController) GetDetailItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := ic.ItemService.GetDetailItem(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(200, gin.H{"item": item})
}

func (ic *ItemController) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}
	
	if err := ic.ItemService.DeleteItem(id); err != nil {
			if err.Error() == "already deleted" {
				c.JSON(400, gin.H{"error": "Item already deleted"})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to delete item", "details": err.Error()})
			return
		}
	c.JSON(200, gin.H{"message": "Item deleted successfully"})
}


func (ic *ItemController) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"error": "Invalid item ID"})
		return
	}

	userID := c.GetUint("user_id")

	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	item.ID = uint(id)
	item.UserID = userID

	if err := ic.ItemService.UpdateItem(&item); err != nil {
		if err.Error() == "item not found or already deleted" {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to update item", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Item updated successfully"})
}
