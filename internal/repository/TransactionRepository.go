package repository

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

// CreatePayment implements [TransactionRepository].
func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
	panic("unimplemented")
}

// FindOrderById implements [TransactionRepository].
func (t *transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	panic("unimplemented")
}

// FindOrders implements [TransactionRepository].
func (t *transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	panic("unimplemented")
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{db: db}
}
