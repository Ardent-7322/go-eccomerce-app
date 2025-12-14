package handlers

import (
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
	handler := TransactionHandler{
		svc:           svc,
		paymentClient: as.Pc,
	}

	secRoute := app.Group("/", as.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	//create payment and collect it

	//1. call user service get cart data to aggregate the total amount and collect payment

	//2. check if payment session active or create a new payment session
	sessionResult, err := h.paymentClient.CreatePayment(100, 123, 456) //amount, userId, orderId
	//3. Store payment session in db to create and validate order

	if err != nil {
		return ctx.Status(400).JSON(err)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":     "payment success!",
		"result":      sessionResult,
		"payment_url": sessionResult.URL,
	})

	payload := struct {
		Message string `json:"message"`
	}{
		Message: "success",
	}
	return ctx.Status(fiber.StatusOK).JSON(payload)
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("success")
}

func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("success")
}
