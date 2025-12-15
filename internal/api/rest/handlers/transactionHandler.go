package handlers

import (
	"errors"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"go-ecommerce-app/pkg/payment"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	svc           *service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) *service.TransactionService {
	return service.NewTransactionService(
		repository.NewTransactionRepository(db),
		auth,
	)
}

func SetupTransactionRoutes(as *rest.RestHandler) {

	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)
	useSvc := service.UserService{
		UserRepo:    repository.NewUserRepository(as.DB),
		CatalogRepo: repository.NewCatalogRepository(as.DB),
		Auth:        as.Auth,
		Config:      as.Config,
	}
	handler := TransactionHandler{
		svc:           svc,
		paymentClient: as.Pc,
		userSvc:       useSvc,
	}

	secRoute := app.Group("/", as.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}
func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	if user.ID == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "unauthenticated",
		})
	}

	// 1. Check active payment
	activePayment, err := h.svc.GetActivePayment(user.ID)
	if err == nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message":     "active payment exists",
			"payment_url": activePayment.PaymentUrl,
		})
	}

	// 2. Get cart total
	_, amount, err := h.userSvc.FindCart(user.ID)
	if err != nil {
		return rest.BadRequestError(ctx, err.Error())
	}

	if amount <= 0 {
		return rest.BadRequestError(ctx, "cart is empty")
	}

	// 3. Generate order reference
	orderRef, err := helper.RandomHandler(8)
	if err != nil {
		return rest.InternalError(ctx, errors.New("error generating order reference"))
	}

	// 4. Create Stripe payment
	sessionResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderRef)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// 5. Store payment
	err = h.svc.StoreCreatedPayment(user.ID, sessionResult, amount, orderRef)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":     "payment created",
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("success")
}

func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("success")
}
