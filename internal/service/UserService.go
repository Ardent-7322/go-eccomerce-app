package service

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"log"
	"time"
)

// UserService handles all business logic related to users
type UserService struct {
	UserRepo repository.UserRepository // DB operations for user
	Auth     helper.Auth               // Auth tools: hashing, token, verify
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

func (s UserService) GetVerificationCode(e domain.User) (string, error) {

	// if user already verified

	if s.isVerifiedUser(e.ID) {
		return "", errors.New("user already verified")
	}

	// generate verifaction code
	code, err := s.Auth.GenerateCode()
	if err != nil {
		return "", err
	}

	// Prepare updated user data
	user := domain.User{
		ID:     e.ID, // ← must include ID
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	// Update in DB
	_, err = s.UserRepo.UpdateUser(e.ID, user)
	if err != nil {
		return "", errors.New("unable to update verification code")
	}
	//send SMS
	// return verification code

	return code, nil
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
func (s UserService) BecomeSeller(id uint, input any) (string, error) {

	return "", nil
}
func (s UserService) FindCart(id uint) ([]interface{}, error) {

	return nil, nil
}
func (s UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {

	return nil, nil
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
