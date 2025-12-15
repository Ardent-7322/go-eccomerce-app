package repository

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindInitialPayment(uId uint) (*domain.Payment, error)
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

// FindPayment implements [TransactionRepository].
func (t *transactionStorage) FindInitialPayment(uId uint) (*domain.Payment, error) {
	var payment *domain.Payment
	err := t.db.First(&payment, "user_id=? AND status=initial", uId).Order("created_at_desc").Error
	return payment, err
}

// CreatePayment implements [TransactionRepository].
func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
	return t.db.Create(payment).Error
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
