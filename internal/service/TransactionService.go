package service

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

type TransactionService struct {
	TransactionRepo repository.TransactionRepository
	Auth            helper.Auth
}

func (s TransactionService) GetOrders(u domain.User) ([]domain.OrderItem, error) {
	orders, err := s.TransactionRepo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s TransactionService) GetOrderDetails(u domain.User, id uint) (dto.SellerOrderDetails, error) {
	orders, err := s.TransactionRepo.FindOrderById(u.ID, id)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}
	return orders, nil
}

func (s TransactionService) GetActivePayment(uId uint) (*domain.Payment, error) {
	return s.TransactionRepo.FindInitialPayment(uId)
}

func (s TransactionService) StoreCreatedPayment(input dto.CreatePaymentRequest) error {
	payment := domain.Payment{
		UserId:       input.UserId,
		Amount:       input.Amount,
		Status:       domain.PaymentStatusInitial,
		PaymentId:    input.PaymentId,
		ClientSecret: input.ClientSecret,
		OrderId:      input.OrderId,
	}
	return s.TransactionRepo.CreatePayment(&payment)
}

func (s TransactionService) UpdatePayment(userId uint, status string, paymentlog string) error {
	p, err := s.GetActivePayment(userId)
	if err != nil {
		return err
	}
	p.Status = domain.PaymentStatus(status)
	p.Response = paymentlog
	return s.TransactionRepo.UpdatePayment(p)
}
func NewTransactionService(r repository.TransactionRepository, auth helper.Auth) *TransactionService {
	return &TransactionService{
		TransactionRepo: r,
		Auth:            auth,
	}
}
