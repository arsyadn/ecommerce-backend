package repositories

import (
	"final-project/models"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (rr *ReportRepository) GetOrders() ([]models.OrderReport, error) {
	var orders []models.OrderReport
	rows, err := rr.DB.Raw(`
		SELECT o.id as order_id, o.user_id, o.status,
			   COALESCE(p.amount, 0) as amount,
			   o.created_at, o.updated_at
		FROM orders o
		LEFT JOIN payments p ON o.id = p.order_id
		ORDER BY o.created_at DESC
	`).Rows()
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.OrderReport
		err := rows.Scan(
			&order.OrderID,
			&order.UserID,
			&order.Status,
			&order.Amount,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, err
}

func (rr *ReportRepository) GetStockMovementsByOrder(orderID int) ([]models.StockMovementReport, error) {
	var stockMovements []models.StockMovementReport
	err := rr.DB.Raw(`
		SELECT id, item_id, reference_order_id as reference_order,
		       quantity, type, created_at
		FROM stockmovement
		WHERE reference_order_id = ?
	`, orderID).Scan(&stockMovements).Error
	return stockMovements, err
}
