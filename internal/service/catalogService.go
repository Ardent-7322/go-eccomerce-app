package service

import (
	"go-ecommerce-app/config"

	"go-ecommerce-app/internal/helper"
	"go-ecommerce-app/internal/repository"
)

// Catalog service handles all business logic related to users
type UserService struct {
	UserRepo repository.UserRepository // DB operations for user
	Auth     helper.Auth               // Auth tools: hashing, token, verify
	Config   config.AppConfig
}
