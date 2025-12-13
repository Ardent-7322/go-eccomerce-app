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

	//cart
	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uId uint, pId uint) (domain.Cart, error)
	CreateCart(c domain.Cart) error
	UpdateCart(c domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(uId uint) error

	//profile
	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

type userRepository struct {
	db *gorm.DB
}

// CreateProfile implements [UserRepository].
func (r *userRepository) CreateProfile(e domain.Address) error {

	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("error on creating profile with address %v", err)
		return errors.New("failed to create profile")
	}

	return nil
}

// UpdateProfile implements [UserRepository].
func (r *userRepository) UpdateProfile(e domain.Address) error {

	err := r.db.Where("user_id=?", e.UserId).Updates(e).Error
	if err != nil {
		log.Printf("error on update profile with address %v", err)
		return errors.New("failed to create profile")
	}
	return nil
}

// CreateCart implements UserRepository.
func (r *userRepository) CreateCart(c domain.Cart) error {
	return r.db.Create(&c).Error
}

// DeleteCartById implements UserRepository.
func (r *userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, id).Error
	return err
}

// DeleteCartItems implements UserRepository.
func (r *userRepository) DeleteCartItems(uId uint) error {
	err := r.db.Where("user_id=?", uId).Delete(&domain.Cart{}).Error
	return err
}

// FindCartItem implements UserRepository.
func (r *userRepository) FindCartItem(uId uint, pId uint) (domain.Cart, error) {
	cartItem := domain.Cart{}
	err := r.db.Where("user_id=? AND product_id=?", uId).Find(&cartItem).Error
	return cartItem, err

}

// FindCartItems implements UserRepository.
func (r *userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id=?", uId).Find(&carts).Error
	return carts, err

}

// UpdateCart implements UserRepository.
func (r *userRepository) UpdateCart(c domain.Cart) error {

	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id=?", c.ID).Updates(c).Error
	return err
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
	err := r.db.Preload("Address").First(&user, "email=?", email).Error

	if err != nil {
		log.Printf("find user error %v", err)
		return domain.User{}, errors.New("user does not exist")
	}
	return user, nil
}

// ✔ matches interface: FindUserById(id uint)
func (r *userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.Preload("Address").First(&user, id).Error

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
