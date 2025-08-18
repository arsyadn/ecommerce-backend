package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	paymentController := controllers.NewPaymentController(db)

	user := router.Group("/")
	user.Use(middleware.AuthMiddleware())
	user.Use(middleware.RoleUserMiddleware())

	admin := router.Group("/")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleAdminMiddleware())
	
	{
		user.POST("/payments", paymentController.UploadPayment)
		admin.GET("/payments", paymentController.GetPayments)
		admin.GET("/payments/:id", paymentController.GetPaymentByID)
		admin.PUT("/payments", paymentController.AdminUpdatePayment)
	}
}