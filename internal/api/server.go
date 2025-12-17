package api

import (
	"log"
	"os"

	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/pkg/payment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(cfg config.AppConfig) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return rest.SuccessResponse(c, "I am Healty", &fiber.Map{
			"status": "ok with 200 status code",
		})
	})

	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Printf("database connection failed: %v", err)
		return
	}
	log.Println("Database connected")

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	); err != nil {
		log.Printf("migration failed: %v", err)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	auth := helper.SetupAuth(cfg.AppSecret)
	paymentClient := payment.NewPaymentClient(cfg.StripeSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: cfg,
		Pc:     paymentClient,
	}

	setupRoutes(rh)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)

	if err := app.Listen(":" + port); err != nil {
		log.Printf("server stopped: %v", err)
	}
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupCatalogRoutes(rh)
	handlers.SetupUserRoutes(rh)
	handlers.SetupTransactionRoutes(rh)
}
