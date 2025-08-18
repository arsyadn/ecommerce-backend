package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupItemsRoutes(router *gin.RouterGroup, db *gorm.DB) {
	itemsController := controllers.NewItemController(db)

	admin := router.Group("/items")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleAdminMiddleware())
	
	protected := router.Group("/items")
	protected.Use(middleware.AuthMiddleware())

	{
		admin.POST("", itemsController.CreateItem) //create item
		protected.GET("", itemsController.GetAllItems)
		protected.GET("/:id", itemsController.GetDetailItem)
		admin.PUT("/:id", itemsController.UpdateItem)
		admin.DELETE("/:id", itemsController.DeleteItem)
	}
}