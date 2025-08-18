package routes

import (
	"final-project/controllers"
	"final-project/middleware"
	"final-project/repositories"
	"final-project/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReportRoutes(router *gin.RouterGroup, db *gorm.DB) {
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportController := controllers.NewReportController(reportService)

	admin := router.Group("/")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RoleAdminMiddleware())
	admin.GET("/admin/reports", reportController.GetAdminReport)
}
