package handlers

import (
	"chatapp/pkg/models"
	"chatapp/pkg/util"
	"chatapp/services/user"
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

type (
	// AuthHandlerOptions represents the options required to set up the auth handler
	AuthHandlerOptions struct {
		UserService user.Service
	}

	// authHandler handles user auth
	authHandler struct {
		userService user.Service
	}
)

// Register adds and returns the new user created
func (h *authHandler) Register(c *fiber.Ctx) error {
	var u *models.User

	if err := c.BodyParser(&u); err != nil {
		return serverError(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := u.ValidateRegisterRequest(); err != nil {
		return validationError(c, err)
	}

	ctx := context.Background()
	now := time.Now().Local()

	u.CreatedAt = now
	u.UpdatedAt = now

	exists, err := h.userService.CheckIfExists(ctx, "username", u.Username)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	if exists {
		return validationDuplicateError(c, fiber.Map{
			"username": "has already been taken",
		})
	}

	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	u.Password = hashedPassword

	newUser, err := h.userService.Create(ctx, u)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(newUser)
}

// Login attempts to log in a user using the provided credentials
func (h *authHandler) Login(c *fiber.Ctx) error {
	return nil
}

// AuthHandler is an interface for the user authentication
type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(opts AuthHandlerOptions) AuthHandler {
	return &authHandler{
		userService: opts.UserService,
	}
}
