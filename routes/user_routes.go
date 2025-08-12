package routes

import (
	"final-project/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	protected := router.Group("/")
	{
		protected.POST("/register", userController.Register)
		protected.POST("/login", userController.Login)
	}
}