package services

import (
	"final-project/models"
	"final-project/repositories"
)

type ReportService struct {
	ReportRepo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{ReportRepo: repo}
}

func (rs *ReportService) GetAdminReport() ([]models.OrderReport, error) {
	orders, err := rs.ReportRepo.GetOrders()
	if err != nil {
		return nil, err
	}

	// attach stock movements per order
	for i := range orders {
		stockMovements, err := rs.ReportRepo.GetStockMovementsByOrder(orders[i].OrderID)
		if err != nil {
			return nil, err
		}
		orders[i].StockMovement = stockMovements
	}

	return orders, nil
}
