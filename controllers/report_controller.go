package controllers

import (
	"final-project/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportService *services.ReportService
}

func NewReportController(service *services.ReportService) *ReportController {
	return &ReportController{ReportService: service}
}

func (rc *ReportController) GetAdminReport(c *gin.Context) {
	reports, err := rc.ReportService.GetAdminReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reports)
}
