package services

import (
	"final-project/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		DB: db,
	}
}

func (os *OrderService) GetUserRole(userID uint) (string, error) {
	var role string
	query := `SELECT role FROM users WHERE id = ?`
	if err := os.DB.Raw(query, userID).Scan(&role).Error; err != nil {
		return "", err
	}
	return role, nil
}

func (os *OrderService) CreateOrder(order *models.Order) error {
	// Set waktu order
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	if order.Status == "" {
		order.Status = "pending"
	}

	return os.DB.Transaction(func(tx *gorm.DB) error {
		
		result := tx.Exec("INSERT INTO orders (created_at, updated_at, status, user_id) VALUES (?, ?, ?, ?)",
			order.CreatedAt, order.UpdatedAt, order.Status, order.UserID)
		if result.Error != nil {
			return result.Error
		}
		
		if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&order.ID).Error; err != nil {
			return err
		}

		for i := range order.Details {
			var count int64
			if err := tx.Raw("SELECT COUNT(*) FROM items WHERE id = ? AND id IN (SELECT id FROM items WHERE id = ?)", 
				order.Details[i].ItemID, order.Details[i].ItemID).Scan(&count).Error; err != nil {
				return  fmt.Errorf("item with ID %d does not exist in the system", order.Details[i].ItemID)
			}
			if count == 0 {
				return fmt.Errorf("item with ID %d not found", order.Details[i].ItemID)
			}
			
			var stock int
			if err := tx.Raw("SELECT stock FROM items WHERE id = ?", order.Details[i].ItemID).Scan(&stock).Error; err != nil {
				return fmt.Errorf("failed to check stock for item with ID %d: %v", order.Details[i].ItemID, err)
			}
			if order.Details[i].Quantity > stock {
				return fmt.Errorf("insufficient stock for item with ID %d: requested %d, available %d", 
					order.Details[i].ItemID, order.Details[i].Quantity, stock)
			}

			var price float64
			if err := tx.Raw("SELECT price FROM items WHERE id = ?", 
				order.Details[i].ItemID).Scan(&price).Error; err != nil {
				return fmt.Errorf("failed to get price for item with ID %d: %v", order.Details[i].ItemID, err)
			}
			fmt.Println("Price at order for item", order.Details[i].ItemID, "is", price)
			order.Details[i].Price = price
		}
		fmt.Println("Order details:", order.Details)
		for i := range order.Details {
			order.Details[i].OrderID = order.ID
			order.Details[i].CreatedAt = time.Now()
			fmt.Println("Inserting order detail for item", order.Details[i].ItemID, "with price at", order.Details[i].PriceAtOrder)


			result := tx.Exec("INSERT INTO orderdetails (order_id, item_id, quantity, price, price_at_order, created_at) VALUES (?, ?, ?, ?, ?, ?)",
				order.Details[i].OrderID, order.Details[i].ItemID, order.Details[i].Quantity, 
				order.Details[i].Price, order.Details[i].PriceAtOrder, order.Details[i].CreatedAt)
			
			if result.Error != nil {
				return fmt.Errorf("failed to insert order detail for item here %d: %v", order.Details[i].ItemID, result.Error)
			}

			if result.RowsAffected == 0 {
				return fmt.Errorf("failed to insert order detail for item %d", order.Details[i].ItemID)
			}
		}

		return nil
	})
}