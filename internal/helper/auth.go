package helper

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func NewAuth() Auth {
	return Auth{
		Secret: os.Getenv("APP_SECRET"),
	}
}

func SetupAuth(secret string) Auth {
	return Auth{
		Secret: secret,
	}
}

/* ================= PASSWORD ================= */

func (a Auth) CreateHashedPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password length should be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password hash failed")
	}

	return string(hash), nil
}

func (a Auth) VerifyPassword(plain string, hashed string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return errors.New("password does not match")
	}
	return nil
}

/* ================= JWT ================= */

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("invalid inputs for token generation")
	}

	if a.Secret == "" {
		return "", errors.New("jwt secret is missing")
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

func (a Auth) VerifyToken(authHeader string) (domain.User, error) {
	if authHeader == "" {
		return domain.User{}, errors.New("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return domain.User{}, errors.New("invalid authorization header format")
	}

	tokenStr := parts[1]

	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil || !parsedToken.Valid {
		return domain.User{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return domain.User{}, errors.New("token expired")
	}

	user := domain.User{
		ID:       uint(claims["user_id"].(float64)),
		Email:    claims["email"].(string),
		UserType: claims["role"].(string),
	}

	return user, nil
}

/* ================= MIDDLEWARE ================= */

func (a Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	user, err := a.VerifyToken(authHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	user, err := a.VerifyToken(authHeader)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err.Error(),
		})
	}

	if user.UserType != domain.SELLER {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  "please join seller program to manage products",
		})
	}

	ctx.Locals("user", user)
	return ctx.Next()
}

/* ================= CONTEXT ================= */

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {
	user, ok := ctx.Locals("user").(domain.User)
	if !ok {
		return domain.User{}
	}
	return user
}

/* ================= OTP ================= */

func (a Auth) GenerateCode() (string, error) {
	return RandomHandler(6)
}
