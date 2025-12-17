package repository

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindInitialPayment(uId uint) (*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

// UpdatePayment implements [TransactionRepository].
func (t *transactionStorage) UpdatePayment(payment *domain.Payment) error {
	return t.db.Save(payment).Error
}

// FindPayment implements [TransactionRepository].
func (t *transactionStorage) FindInitialPayment(uId uint) (*domain.Payment, error) {

	var payment domain.Payment

	err := t.db.
		Where("user_id = ? AND status = ?", uId, domain.PaymentStatusInitial).
		Order("created_at DESC").
		First(&payment).
		Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// CreatePayment implements [TransactionRepository].
func (t *transactionStorage) CreatePayment(payment *domain.Payment) error {
	return t.db.Create(payment).Error
}

// FindOrderById implements [TransactionRepository].
func (t transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	//TODO implement me
	panic("implement me")
}

// FindOrders implements [TransactionRepository].

func (t transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	//TODO implement me
	panic("implement me")
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{db: db}
}
