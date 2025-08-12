package services

import (
	"final-project/models"

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
	query := `INSERT INTO items (name, description, price, stock, users_id) 
			  VALUES (?, ?, ?, ?, ?)`
	if err := is.DB.Exec(query, 
		item.Name, 
		item.Description, 
		item.Price, 
		item.Stock, 
		item.UserID,
		).Error; err != nil {
		return err
	}
	return nil
}

func (is *ItemService) GetAllItems(page, limit int) ([]models.ItemResponse, error) {
	var items []models.ItemResponse
	query := `SELECT id, name, price, stock FROM items LIMIT ? OFFSET ?`
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
	query := `SELECT id, name, description, price, stock FROM items WHERE id = ?`
	if err := is.DB.Raw(query, id).Scan(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (is *ItemService) DeleteItem(id int) error {
	query := `DELETE FROM items WHERE id = ?`
	if err := is.DB.Exec(query, id).Error; err != nil {
		return err
	}
	return nil
}

func (is *ItemService) UpdateItem(item *models.Item) error {
	query := `UPDATE items SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?`
	if err := is.DB.Exec(query, item.Name, item.Description, item.Price, item.Stock, item.ID).Error; err != nil {
		return err
	}
	return nil
}