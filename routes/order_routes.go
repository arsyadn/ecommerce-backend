package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func SetupOrderRoutes(router *gin.RouterGroup, db *gorm.DB) {
	orderController := controllers.NewOrderController(db)

	protected := router.Group("/orders")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/create", orderController.CreateOrder)
	}
}