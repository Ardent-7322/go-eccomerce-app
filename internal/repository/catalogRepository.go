package repository

import (
	"errors"
	"go-ecommerce-app/internal/domain"
	"log"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(e *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id int) (*domain.Category, error)
	EditCategory(e *domain.Category) (*domain.Category, error)
	DeleteCategory(id int) error
}

type catalogRepository struct {
	db *gorm.DB
}

func (c catalogRepository) CreateCategory(e *domain.Category) error {
	err := c.db.Create(&e).Error
	if err != nil {
		log.Printf("db_err: %v", err)
		return errors.New("create category failed")
	}
	return nil
}

func (c catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category

	err := c.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}
	return categories, nil
}
func (c catalogRepository) FindCategoryById(id int) (*domain.Category, error) {
	var category domain.Category
	err := c.db.First(&category, id).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("category does not exist")
	}
	return &category, nil
}
func (c catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {
	err := c.db.Save(&e).Error

	if err != nil {
		log.Printf("db_err: %v", err)
		return nil, errors.New("failed to update category") // ✅ FIXED
	}

	return e, nil
}

func (c catalogRepository) DeleteCategory(id int) error {

	err := c.db.Delete(&domain.Category{}, id).Error
	if err != nil {
		log.Printf("db_err: %v", err)
		return errors.New("failed to delete category")
	}
	return nil
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
