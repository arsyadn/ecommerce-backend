package models

import "time"

type Order struct {
	ID 	  uint           `gorm:"primaryKey" json:"id"`
	UserID   uint           `json:"user_id"`
	Status    string         `gorm:"type:enum('pending','processed','failed','success');default:'pending'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Details   []OrderDetail  `gorm:"foreignKey:OrderID" json:"details"`
}

type OrderDetail struct {
	OrderID     uint      `json:"order_id"`
	ItemID      uint      `json:"item_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	PriceAtOrder float64  `json:"price_at_order"`
	CreatedAt   time.Time `json:"created_at"`
}