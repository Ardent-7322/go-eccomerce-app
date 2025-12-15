package domain

import "time"

type BankAccount struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserId      uint      `json:"user_id" gorm:"index;not null"`
	BankAccount string    `json:"bank_account_number" gorm:"not null"`
	SwiftCode   string    `json:"swift_code" gorm:"not null"`
	PaymentType string    `json:"payment_type" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
