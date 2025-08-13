package models

import "time"

type PaymentRequest struct {
	OrderID  uint    `json:"order_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Evidence []byte `json:"evidence"`
}

type Payment struct {
	ID        int       `json:"id"`
	OrderID   uint      `json:"order_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Evidence  []byte    `json:"evidence"`
	CreatedAt time.Time `json:"created_at"`
}