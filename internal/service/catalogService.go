package service

import (
	"errors"
	"go-ecommerce-app/config"

	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

// Catalog service handles all business logic related to users
type UserService struct {
	UserRepo repository.UserRepository // DB operations for user
	Auth     helper.Auth               // Auth tools: hashing, token, verify
	Config   config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {

	err := s.CatalogRepo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryRequest) (*domain.Category, error) {

	existCat, err := s.CatalogRepo.FindCategoryById(id)

	if err != nil {
		return nil, errors.New("category does not exist")
	}

	if len(input.Name) > 0 {
		existCat.Name = input.Name
	}

	if input.ParentId > 0 {
		existCat.ParentId = input.ParentId
	}

	if len(input.ImageUrl) > 0 {
		existCat.ImageUrl = input.ImageUrl
	}

	if input.DisplayOrder > 0 {
		existCat.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.CatalogRepo.EditCategory(existCat)

	return updatedCat, err

}

func (s CatalogService) DeleteCategory(id int) error {

	err := s.CatalogRepo.DeleteCategory(id)
	if err != nil {
		// log the error
		return errors.New("category does not exist to delete")
	}
	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {

	categories, err := s.CatalogRepo.FindCategories()
	if err != nil {
		return nil, errors.New("categories do not exist")
	}

	return categories, nil
}

func (s CatalogService) GetCategory(id int) (*domain.Category, error) {
	cat, err := s.CatalogRepo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}
	return cat, nil
}
