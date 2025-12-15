package service

import (
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"

	"github.com/stripe/stripe-go/v78"
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

func (s TransactionService) StoreCreatedPayment(uId uint, ps *stripe.CheckoutSession, amount float64, orderId string) error {
	payment := domain.Payment{
		UserId:     uId,
		Amount:     amount,
		Status:     domain.PaymentStatusInitial,
		PaymentUrl: ps.URL,
		PaymentId:  ps.ID,
		OrderId:    orderId,
	}
	return s.TransactionRepo.CreatePayment(&payment)
}

func NewTransactionService(r repository.TransactionRepository, auth helper.Auth) *TransactionService {
	return &TransactionService{
		TransactionRepo: r,
		Auth:            auth,
	}
}
