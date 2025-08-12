package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupItemsRoutes(router *gin.RouterGroup, db *gorm.DB) {
	itemsController := controllers.NewItemController(db)

	protected := router.Group("/items")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/create", itemsController.CreateItem)
		protected.GET("", itemsController.GetAllItems)
		protected.GET("/:id", itemsController.GetDetailItem)
		protected.PUT("/:id", itemsController.UpdateItem)
		protected.DELETE("/:id", itemsController.DeleteItem)
	}
}