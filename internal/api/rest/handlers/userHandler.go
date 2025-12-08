package handlers

import (
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/dto"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// svc UserService
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	//create an instance of user service & inject to handler
	svc := service.UserService{
		UserRepo: repository.NewUserRepository(rh.DB),
		Auth:     rh.Auth,
	}
	handler := UserHandler{
		svc: svc,
	}
	//Grouping kardenge
	pubRoutes := app.Group("/users")
	//Public endpoints
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.login)

	//Private routes ko grouping kardenge and can be accessible only by authorization
	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)
	//private endpoints
	pvtRoutes.Get("/verify", handler.GetverificationCode)
	pvtRoutes.Post("/verify", handler.verify)
	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)

	pvtRoutes.Post("/cart", handler.AddtoCart)
	pvtRoutes.Get("/cart", handler.GetCart)
	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrders)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignup{}
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
			"error":   err.Error(), // optional
		})
	}

	token, err := h.svc.Signup(user)
	if err != nil {
		log.Printf("Signup error: %v\n", err) // <--- THIS IS IMPORTANT

		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error on signup",
			"error":   err.Error(), // add this for now to see real cause
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "signup success",
		"token":   token,
	})
}
func (h *UserHandler) login(ctx *fiber.Ctx) error {

	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Please provide correct user id password",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "login",
		"token":   token,
	})
}
func (h *UserHandler) GetverificationCode(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetverificationCode",
	})
}
func (h *UserHandler) verify(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "verify",
	})
}
func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "CreateProfile",
	})
}
func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetProfile",
		"user":    user,
	})
}
func (h *UserHandler) AddtoCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "AddtoCart",
	})
}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetCart",
	})
}
func (h *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "CreateOrder",
	})
}
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetOrders",
	})
}
func (h *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "BecomeSeller",
	})
}
