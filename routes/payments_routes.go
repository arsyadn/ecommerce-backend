package routes

import (
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	paymentController := controllers.NewPaymentController(db)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/payments", paymentController.UploadPayment)
		protected.GET("/payments", paymentController.GetPayments)
	}
}