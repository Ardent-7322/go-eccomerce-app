package handlers

import (
	"fmt"
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
		Config:   rh.Config,
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
	user := h.svc.Auth.GetCurrentUser(ctx)

	err := h.svc.GetVerificationCode(user)
	if err != nil {
		fmt.Println("GetVerificationCode error:", err) // 🔍 DEBUG
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(), // temporarily return real reason
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "GetverificationCode",
	})
}
func (h *UserHandler) verify(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	//request
	var req dto.VerificationCodeInput

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid input",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)

	if err != nil {
		log.Printf("%v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "verified successfully",
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

	req := dto.CreateCarRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid product and qty",
		})
	}
	user := h.svc.Auth.GetCurrentUser(ctx)

	//call user service and perform create cart
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "cart created successfully ", cartItems)

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

	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "failed to become seller",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "become seller",
		"token":   token,
	})
}
