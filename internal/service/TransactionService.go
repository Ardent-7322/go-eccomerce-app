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

func NewTransactionService(r repository.TransactionRepository, auth helper.Auth) *TransactionService {
	return &TransactionService{
		TransactionRepo: r,
		Auth:            auth,
	}
}
