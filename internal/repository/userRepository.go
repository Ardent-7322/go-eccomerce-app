package repository

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	CreateUser(usr domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)

	//more function will come as we progress

	CreateBankAccount(e domain.BankAccount) error
}

type userRepository struct {
	db *gorm.DB
}

func (r userRepository) CreateBankAccount(e domain.BankAccount) error {
	return r.db.Create(&e).Error
}
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// ✔ pointer receiver
func (r *userRepository) CreateUser(usr domain.User) (domain.User, error) {
	err := r.db.Create(&usr).Error
	if err != nil {
		log.Printf("create user error %v", err)
		// 👇 yahan generic error mat bhejo, real error wrap karo
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return usr, nil
}

// ✔ matches interface: FindUser(email string)
func (r *userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		log.Printf("find user error %v", err)
		return domain.User{}, errors.New("user does not exist")
	}
	return user, nil
}

// ✔ matches interface: FindUserById(id uint)
func (r *userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error

	if err != nil {
		log.Printf("find user error %v", err)
		return domain.User{}, errors.New("user does not exist")
	}
	return user, nil
}

// ✔ matches interface: UpdatedUser(id uint, u domain.User)
func (r *userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return domain.User{}, err
	}

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(u).Error
	if err != nil {
		log.Printf("Update user error %v", err)
		return domain.User{}, errors.New("failed to update user")
	}

	return user, nil
}
