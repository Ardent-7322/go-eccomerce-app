package helper

import (
	"errors"
	"fmt"
	"go-ecommerce-app/internal/domain"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"strings"
	"time"

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
func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(password string) (string, error) {

	if len(password) < 6 {
		return "", errors.New("password length should be atleast 6 characters long")
	}

	hashP, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		// log actual error and report to logging tool
		return "", errors.New("password hash failed")
	}

	return string(hashP), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}

	if a.Secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
		"iat":     time.Now().Unix(),
	}

	// HS256 use karo, ES256 nahi
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//  string nahi, []byte(secret) pass karo
	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		// yahan original err swallow mat karo
		return "", fmt.Errorf("unable to sign the token: %w", err)
	}

	return tokenStr, nil
}

func (a Auth) VerifyPassword(plain_Password string, hashed_Password string) error {
	if len(plain_Password) < 6 {
		return errors.New("password length should be atleast 6 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed_Password), []byte(plain_Password))

	if err != nil {
		return errors.New("Password does not match")
	}
	return nil

}

func (a Auth) VerifyToken(authHeader string) (domain.User, error) {
	// Expected format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return domain.User{}, errors.New("invalid authorization header")
	}

	tokenStr := parts[1]

	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure token is signed with HMAC (HS256/384/512)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		return domain.User{}, errors.New("invalid token")
	}

	// Validate claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return domain.User{}, errors.New("invalid token claims")
	}

	// Expiry check
	exp, ok := claims["exp"].(float64)
	if !ok {
		return domain.User{}, errors.New("invalid exp claim")
	}
	if time.Now().Unix() > int64(exp) {
		return domain.User{}, errors.New("token is expired")
	}

	// Build user from claims
	user := domain.User{
		ID:       uint(claims["user_id"].(float64)),
		Email:    claims["email"].(string),
		UserType: claims["role"].(string),
	}

	return user, nil
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)

	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "autherization failed",
			"reason":  err,
		})
	}
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) domain.User {
	user := ctx.Locals("user")
	return user.(domain.User)
}

func (a Auth) GenerateCode() (string, error) {
	return RandomHandler(6)
}

func (a Auth) AuthorizeSeller(ctx *fiber.Ctx) error {

	authHeader := ctx.Get("Authorization")
	user, err := a.VerifyToken(authHeader)

	if err != nil {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})

	} else if user.ID > 0 && user.UserType == domain.SELLER {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  errors.New("please join seller program to manage products"),
		})
	}
}
