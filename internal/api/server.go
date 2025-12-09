package api

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/config"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/api/rest/handlers"
	"go-ecommerce-app/internal/domain"
	"go-ecommerce-app/internal/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func StartServer(cfg config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error %v\n", err)
	}
	log.Println("Database connected")

	db.AutoMigrate(&domain.User{})

	auth := helper.SetupAuth(cfg.AppSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: cfg, //  use cfg here
	}
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
}
