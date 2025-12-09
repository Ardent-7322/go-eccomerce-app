package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	// svc UserService
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	//create an instance of user service & inject to handler
	svc := service.CatalogService{
		CatalogRepo: repository.NewCatalogRepository(rh.DB),
		Auth:        rh.Auth,
		Config:      rh.Config,
	}
	handler := CatalogHandler{
		svc: svc,
	}

	// public
	// listing products and categories

	app.Get("/products")
	app.Get("/products/:id")
	app.Get("/categories")
	app.Get("/catagories/:id")

	// private
	// manage product and categories
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	// Categories
	selRoutes.Post("/categories", handler.CreateCategories)
	selRoutes.Patch("/categories/:id", handler.EditCategory)
	selRoutes.Delete("/categories/:id", handler.DeleteCategory)

	// Products
	selRoutes.Post("/products", handler.CreateProducts)
	selRoutes.Get("/products", handler.GetProducts)
	selRoutes.Get("/products/:id", handler.GetProduct)
	selRoutes.Patch("/products/:id", handler.EditProduct)
	selRoutes.Patch("/products/:id", handler.UpdateStock) // update stock
	selRoutes.Delete("/products/:id", handler.DeleteCategory)

}

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	log.Printf("current user %v", user.ID)
	return rest.SuccessResponse(ctx, "Create categories endpoint", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "edit category endpoint", nil)

}
func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Delete category endpoint", nil)
}
func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Create category endpoint", nil)
}
func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Edit product endpoint", nil)
}
func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get products by ID", nil)
}
func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get products by ID", nil)
}
func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "category endpoint", nil)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "category endpoint", nil)
}
