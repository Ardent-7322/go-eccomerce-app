package api

import (
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(cfg config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error %v\n", err)
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
	)
	if err != nil {
		log.Fatal("error on running migration %v", err.Error())
	}
	log.Println("migration was successfull")

	//cors configiration

	// cors configuration - expanded for both localhost and 127.0.0.1 (dev only)
	c := cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3030, http://127.0.0.1:3030",
		AllowHeaders:     "Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	})
	app.Use(c)

	auth := helper.SetupAuth(cfg.AppSecret)

	rh := &rest.RestHandler{
		App: app,
		DB:  db, Auth: auth,
		Config: cfg, //  use cfg here
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
			// return immediately for preflight
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	setupRoutes(rh)

	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}

func setupRoutes(rh *rest.RestHandler) {

	//user handlers
	handlers.SetupUserRoutes(rh)

	//transaction handlers

	//catalog

	handlers.SetupCatalogRoutes(rh)
}
