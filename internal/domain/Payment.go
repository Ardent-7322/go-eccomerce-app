package domain

import "time"

type Payment struct {
	ID            uint          `gorm:"PrimaryKey" json:"id"`
	UserId        uint          `json:"user_id"`
	CaptureMethod string        `json:"capture_method"`
	Amount        float64       `json:"amount"`
	TransactionId uint          `json:"transaction_id"`
	OrderId       string        `json:"order_id"`
	CustomerId    string        `json:"customer_id"`                   // stripe customer id
	PaymentId     string        `json:"payment_id"`                    // payment id
	Status        PaymentStatus `json:"status" gorm:"default:initial"` // initial, success, failed
	Response      string        `json:"response"`
	PaymentUrl    string        `json:"payment_url"`
	CreatedAt     time.Time     `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "initial"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)
