package services

import (
	"database/sql"
	"final-project/models"
	"fmt"

	"gorm.io/gorm"
)

type ItemService struct {
	DB *gorm.DB
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{
		DB: db,
	}
}

func (is *ItemService) GetUserRole(userID uint) (string, error) {
	var role string
	query := `SELECT role FROM users WHERE id = ?`
	if err := is.DB.Raw(query, userID).Scan(&role).Error; err != nil {
		return "", err
	}
	return role, nil
}

func (is *ItemService) CreateItem(item *models.Item) error {
	// Start transaction
	tx := is.DB.Begin()

	// Insert into items table
	itemQuery := `INSERT INTO items (name, description, price, stock, users_id) 
				  VALUES (?, ?, ?, ?, ?)`
	result := tx.Exec(itemQuery, 
		item.Name, 
		item.Description, 
		item.Price, 
		item.Stock, 
		item.UserID)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// Get the last inserted item ID
	var itemID uint
	if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&itemID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Insert into stockmovement table
	stockQuery := `INSERT INTO stockmovement (item_id, quantity, type, created_at, created_by) 
				   VALUES (?, ?, 'in', CURRENT_TIMESTAMP, ?)`
	if err := tx.Exec(stockQuery, itemID, item.Stock, item.UserID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

func (is *ItemService) GetAllItems(page, limit int) ([]models.ItemResponse, error) {
	var items []models.ItemResponse
	query := `SELECT id, name, price, stock FROM items WHERE deleted_at IS NULL LIMIT ? OFFSET ?`
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if page <= 0 {
		page = 1 // Default page
	}
	offset := (page - 1) * limit
	if err := is.DB.Raw(query, limit, offset).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (is *ItemService) GetDetailItem(id int) (*models.ItemDetailReponse, error) {
	var item models.ItemDetailReponse
	query := `SELECT id, name, description, price, stock FROM items WHERE id = ? AND deleted_at IS NULL`
	if err := is.DB.Raw(query, id).Scan(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (is *ItemService) DeleteItem(id int) error {
	// Check if item exists
	var count int64
	existsQuery := `SELECT COUNT(*) FROM items WHERE id = ?`
	if err := is.DB.Raw(existsQuery, id).Scan(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("item not found")
	}

	// Check if item already deleted
	var deletedAt sql.NullTime
	checkQuery := `SELECT deleted_at FROM items WHERE id = ?`
	if err := is.DB.Raw(checkQuery, id).Scan(&deletedAt).Error; err != nil {
		return err
	}

	if deletedAt.Valid {
		return fmt.Errorf("already deleted")
	}

	// Soft delete the item
	deleteQuery := `UPDATE items SET deleted_at = NOW() WHERE id = ?`
	if err := is.DB.Exec(deleteQuery, id).Error; err != nil {
		return err
	}

	return nil
}


func (is *ItemService) UpdateItem(item *models.Item) error {
	query := `UPDATE items SET name = ?, description = ?, price = ?, stock = ? WHERE id = ? AND deleted_at IS NULL`
	if err := is.DB.Exec(query, item.Name, item.Description, item.Price, item.Stock, item.ID).Error; err != nil {
		return err
	}
	return nil
}