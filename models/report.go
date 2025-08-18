package models

import "time"

type StockMovementReport struct {
	ID             int       `json:"id"`
	ItemID         int       `json:"item_id"`
	ReferenceOrder int       `json:"reference_order"`
	Quantity       int       `json:"quantity"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
}

type OrderReport struct {
	OrderID       int                   `json:"order_id"`
	UserID        int                   `json:"user_id"`
	Status        string                `json:"status"`
	Amount        float64               `json:"amount"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	StockMovement []StockMovementReport `json:"stock_movements"`
}
