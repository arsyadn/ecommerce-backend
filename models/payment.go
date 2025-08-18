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
	EvidencePath  string    `json:"evidence_path"`
	Evidence []byte `json:"evidence"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentResponse struct {
	ID        int       `json:"id"`
	OrderID   uint      `json:"order_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	EvidencePath  string    `json:"evidence_path"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdatePaymentStatusPayload struct {
	PaymentID int    `json:"payment_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}