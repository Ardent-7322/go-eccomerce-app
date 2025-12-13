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
type CatalogService struct {
	UserRepo    repository.UserRepository
	CatalogRepo repository.CatalogRepository // DB operations for user
	Auth        helper.Auth                  // Auth tools: hashing, token, verify
	Config      config.AppConfig
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

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.CatalogRepo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		UserId:      int(user.ID),
		Stock:       int(input.Stock),
	})
	return err
}

func (s CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {

	existProduct, err := s.CatalogRepo.FindProductById(id)

	// verify product owner
	if existProduct.UserId != int(user.ID) {
		return nil, errors.New("you don't have manage rights of this product")
	}

	if err != nil {
		return nil, errors.New("product does not exist")
	}

	if len(input.Name) > 0 {
		existProduct.Name = input.Name
	}

	if len(input.Description) > 0 {
		existProduct.Description = input.Description
	}

	if input.Price > 0 {
		existProduct.Price = input.Price
	}

	if input.CategoryId > 0 {
		existProduct.CategoryId = input.CategoryId
	}

	updatedProduct, err := s.CatalogRepo.EditProduct(existProduct)

	return updatedProduct, err

}

func (s CatalogService) DeleteProduct(id int, user domain.User) error {
	existProduct, err := s.CatalogRepo.FindProductById(id)
	if err != nil {
		return errors.New("product does not exist")
	}

	// verify product owner
	if existProduct.UserId != int(user.ID) {
		return errors.New("you don't have manage rights of this product")
	}
	err = s.CatalogRepo.DeleteProduct(existProduct)
	if err != nil {
		return errors.New("product cannot delete")
	}

	return nil
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.CatalogRepo.FindProducts()
	if err != nil {
		return nil, errors.New("products do not exist")
	}

	return products, nil
}

func (s CatalogService) GetProductsById(id int) (*domain.Product, error) {
	product, err := s.CatalogRepo.FindProductById(id)
	if err != nil {
		return nil, errors.New("products do not exist")
	}

	return product, nil
}

func (s CatalogService) GetSellerProducts(id int) ([]*domain.Product, error) {
	products, err := s.CatalogRepo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("products do not exist")
	}

	return products, nil
}

func (s CatalogService) UpdatedProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.CatalogRepo.FindProductById(int(e.ID))
	if err != nil {
		return nil, errors.New("product not found")
	}

	// verify product owner
	if product.UserId != e.UserId {
		return nil, errors.New("you don't have manage rights of this product")
	}
	product.Stock = e.Stock
	editProduct, err := s.CatalogRepo.EditProduct(product)
	if err != nil {
		return nil, err
	}
	return editProduct, nil
}
