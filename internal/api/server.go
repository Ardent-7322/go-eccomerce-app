package api

import (
	"log"

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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	log.Println("Database connected")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)
	if err != nil {
		log.Printf("migration error: %v", err)
	}
	log.Println("Migration completed")

	c := cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3030, http://127.0.0.1:3030",
		AllowHeaders:     "Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	})
	app.Use(c)

	app.Get("/", func(c *fiber.Ctx) error {
		return rest.SuccessResponse(c, "Service is running", &fiber.Map{
			"status": "ok",
		})
	})

	auth := helper.SetupAuth(cfg.AppSecret)
	paymentClient := payment.NewPaymentClient(cfg.StripeSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: cfg,
		Pc:     paymentClient,
	}

	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	setupRoutes(rh)

	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupCatalogRoutes(rh)
	handlers.SetupUserRoutes(rh)
	handlers.SetupTransactionRoutes(rh)
}
