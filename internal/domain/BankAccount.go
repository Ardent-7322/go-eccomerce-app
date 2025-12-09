package domain

import "time"

type BankAccount struct {
	ID          uint      `json:"id" gorm:"PrimaryKey"`
	UserId      uint      `json:"user_id"`
	BankAccount uint      `json:"bank_account" gorm:"index;unique;not null"`
	SwiftCode   string    `json:"email" gorm:"index;unique;not null"`
	PaymentType string    `json:"phone"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
