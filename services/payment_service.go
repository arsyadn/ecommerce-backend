package services

import (
	"context"
	"errors"
	"final-project/models"
	"io"
	"time"

	"gorm.io/gorm"
)

type PaymentService struct {
	DB *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{
		DB: db,
	}
}

func (ps *PaymentService) GetUserRole(userID uint) (string, error) {
	var role string
	query := `SELECT role FROM users WHERE id = ?`
	if err := ps.DB.Raw(query, userID).Scan(&role).Error; err != nil {
		return "", err
	}
	return role, nil
}


func (ps *PaymentService) UploadPayment(ctx context.Context, orderID uint, amount float64, evidenceFile io.Reader) error {
    evidenceBytes, err := io.ReadAll(evidenceFile)
    if err != nil {
        return err
    }

    // Check order exists
    var orderCount int64
    if err := ps.DB.WithContext(ctx).
        Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderID).
        Scan(&orderCount).Error; err != nil {
        return err
    }
    if orderCount == 0 {
        return errors.New("order not found")
    }

    // Validate order total price
    var orderPrice float64
    if err := ps.DB.WithContext(ctx).Raw(`
        SELECT COALESCE(SUM(od.price_at_order * od.quantity), 0) 
        FROM orderdetails od 
        WHERE od.order_id = ?`, orderID).Scan(&orderPrice).Error; err != nil {
        return err
    }

    if amount != orderPrice {
        return errors.New("payment amount does not match order total")
    }

    // Create payment record
    payment := models.Payment{
        OrderID:   orderID,
        Status:    "pending",
        CreatedAt: time.Now(),
        Amount:    amount,
        Evidence:  evidenceBytes,
    }

    tx := ps.DB.WithContext(ctx).Begin()
    if tx.Error != nil {
        return tx.Error
    }
    if err := tx.Exec(`
        INSERT INTO payments (order_id, status, created_at, amount, evidence, payment_date) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        payment.OrderID, payment.Status, payment.CreatedAt, payment.Amount, payment.Evidence, payment.CreatedAt).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Exec("UPDATE orders SET status = ? WHERE id = ?", "processed", orderID).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}


func (ps *PaymentService) GetPayments(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	if err := ps.DB.WithContext(ctx).Raw("SELECT * FROM payments").Scan(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (ps *PaymentService) GetPaymentByID(ctx context.Context, paymentID int) (*models.Payment, error) {
	var payment models.Payment
	if err := ps.DB.WithContext(ctx).Raw("SELECT * FROM payments WHERE id = ?", paymentID).Scan(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (ps *PaymentService) AdminUpdatePayment (ctx context.Context, paymentID int, status string, userID int) error {
    var payment models.Payment
    // Get payment
    if err := ps.DB.WithContext(ctx).Raw("SELECT * FROM payments WHERE id = ?", paymentID).Scan(&payment).Error; err != nil {
        return err
    }
    // Validate payment status
    if status != "success" && status != "failed" {
        return errors.New("invalid payment status")
    }

    // Update payment status
    if err := ps.DB.WithContext(ctx).Exec("UPDATE payments SET status = ?, updated_at = ?, updated_by = ? WHERE id = ?", 
        status, time.Now(), userID, paymentID).Error; err != nil {
        return err
    }

    // Update order status based on payment status
    orderStatus := "failed"
    if status == "success" {
        orderStatus = "processed"
    }
    if err := ps.DB.WithContext(ctx).Exec(`
        UPDATE orders o 
        SET o.status = ?
        WHERE o.id = (SELECT p.order_id FROM payments p WHERE p.id = ?)`,
        orderStatus, paymentID).Error; err != nil {
        return err
    }

    if status == "success" {
        // Get order details for the payment
        type OrderDetail struct {
            ItemID   uint
            Quantity int
        }
        var orderDetails []OrderDetail
        if err := ps.DB.WithContext(ctx).Raw(`
            SELECT od.item_id, od.quantity 
            FROM payments p 
            JOIN orderdetails od ON od.order_id = p.order_id 
            WHERE p.id = ?`, paymentID).Scan(&orderDetails).Error; err != nil {
            return err
        }

        // Update stock for each item
        for _, od := range orderDetails {
            // Check if stock is sufficient
            var currentStock int
            if err := ps.DB.WithContext(ctx).Raw(`
            SELECT stock FROM items WHERE id = ?`, od.ItemID).Scan(&currentStock).Error; err != nil {
            return err
            }
            
            if currentStock < od.Quantity {
            return errors.New("insufficient stock for item")
            }

            if err := ps.DB.WithContext(ctx).Exec(`
            UPDATE items 
            SET stock = stock - ? 
            WHERE id = ?`, od.Quantity, od.ItemID).Error; err != nil {
            return err
            }
        }
    }
    return nil
}