package service

import (
	"errors"
	"fmt"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/notification"
	"log"
	"time"
)

// UserService handles all business logic related to users

type UserService struct {
	UserRepo    repository.UserRepository // DB operations for user
	CatalogRepo repository.CatalogRepository
	Auth        helper.Auth // Auth tools: hashing, token, verify
	Config      config.AppConfig
}

// Signup creates a new user and returns a JWT
func (s *UserService) Signup(input dto.UserSignup) (string, error) {

	// Hash the user's password before storing it
	hashedPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	// Save the user in the database
	user, err := s.UserRepo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})
	if err != nil {
		return "", err
	}

	// Generate JWT token for the newly created user
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

// findUserByEmail fetches a user from DB using email (internal helper method)
func (s *UserService) findUserByEmail(email string) (*domain.User, error) {
	user, err := s.UserRepo.FindUser(email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Login verifies credentials and returns a JWT
func (s *UserService) Login(email, password string) (string, error) {

	// Find user by email
	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("user does not exist with the provided email id")
	}

	// Compare hashed password with user input
	if err := s.Auth.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	// Generate JWT token for the user
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.UserRepo.FindUserById(id)

	return err == nil && currentUser.Verified

}

func (s UserService) GetVerificationCode(e domain.User) error {
	// 1) Block already verified users (signup verification flow)
	if s.isVerifiedUser(e.ID) {
		return errors.New("user already verified")
	}

	// 2) Generate verification code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// 3) Update user with code + expiry
	e.Code = code
	e.Expiry = time.Now().Add(30 * time.Minute)

	if _, err := s.UserRepo.UpdateUser(e.ID, e); err != nil {
		return fmt.Errorf("unable to update verification code: %w", err)
	}

	// 4) Re-fetch user to be sure we have fresh data (phone etc.)
	dbUser, err := s.UserRepo.FindUserById(e.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch user after update: %w", err)
	}

	// 5) Send SMS
	notificationClient := notification.NewNotificationClient(s.Config)

	msg := fmt.Sprintf("Your verification code is %v", code)
	if err := notificationClient.SendSMS(dbUser.Phone, msg); err != nil {
		return fmt.Errorf("error sending sms: %w", err)
	}

	return nil
}

func (s UserService) VerifyCode(id uint, code string) error {
	// verify logic here
	if s.isVerifiedUser(id) {
		log.Println("verified...")
		return errors.New("user already verified")
	}

	user, err := s.UserRepo.FindUserById(id)

	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("Verifcation code does not match")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	updateUser := domain.User{
		Verified: true,
	}

	_, err = s.UserRepo.UpdateUser(id, updateUser)

	if err != nil {
		return errors.New("unable to verify user")
	}
	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	return nil, nil
}
func (s UserService) UpdateProfile(id uint, input any) error {
	return nil
}
func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {

	//find existing user
	user, _ := s.UserRepo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("You are already a seller.")
	}

	// update user
	seller, err := s.UserRepo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	// generating token
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// create bank account information

	err = s.UserRepo.CreateBankAccount(domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	})

	return token, err
}
func (s UserService) FindCart(id uint) ([]any, error) {

	return nil, nil
}
func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	// check if the cart is Exist

	cart, _ := s.UserRepo.FindCartItem(u.ID, input.ProductId)

	if cart.ID > 0 {
		if input.ProductId == 0 {
			return nil, errors.New("Please provide a valid product id")
		}
		// -> delete the cart item
		if input.Qty < 1 {
			err := s.UserRepo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error on deleting cart item %v", err)
				return nil, errors.New("error on deleting cart item")
			}

		} else {
			// => update the cart item
			cart.Qty = input.Qty
			err := s.UserRepo.UpdateCart(cart)
			if err != nil {
				//log error
				return nil, errors.New("error on updating cart item")
			}
		}

	} else {
		// check if product exists
		product, _ := s.CatalogRepo.FindProductById(int(input.ProductId))
		if product.ID > 0 {
			return nil, errors.New("product not found")
		}

		// create cart
		err := s.UserRepo.CreateCart(domain.Cart{
			UserId:    u.ID,
			ProductId: input.ProductId,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       input.Qty,
			Price:     product.Price,
			SellerId:  uint(product.UserId),
		})
		if err != nil {
			return nil, errors.New("error on creating cart item")
		}
	}

	return s.UserRepo.FindCartItems(u.ID)
}
func (s UserService) CreateOrder(u domain.User) (int, error) {

	return 0, nil
}
func (s UserService) GetOrders(u domain.User) ([]interface{}, error) {

	return nil, nil
}
func (s UserService) GetOrderById(id uint, uId uint) ([]interface{}, error) {

	return nil, nil
}
