package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func SetupOrderRoutes(router *gin.RouterGroup, db *gorm.DB) {
	orderController := controllers.NewOrderController(db)

	user := router.Group("/orders")
	user.Use(middleware.AuthMiddleware())
	user.Use(middleware.RoleUserMiddleware())

	{
		user.POST("", orderController.CreateOrder)
	}
}