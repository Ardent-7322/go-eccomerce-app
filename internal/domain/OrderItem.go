package domain

import "time"

type OrderItem struct {
	ID        uint      `json:"id" gorm:"PrimaryKey"`
	OrderId   uint      `json:"user_id"`
	ProductId uint      `json:""`
	Name      string    `json:"amount" `
	ImageUrl  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
