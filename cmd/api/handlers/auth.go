package handlers

import (
	"chatapp/pkg/models"
	"chatapp/pkg/util"
	"chatapp/services/user"
	"errors"
	"github.com/gofiber/fiber/v2"
	"time"
)

var (
	errInvalidCredentials = "Invalid username or password provided."
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

	ctx := c.Context()
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

	return successResponse(c, fiber.StatusCreated, fiber.Map{
		"user": models.User{
			ID:        newUser.ID,
			Username:  newUser.Username,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		},
	})
}

// Login attempts to log in a user using the provided credentials
func (h *authHandler) Login(c *fiber.Ctx) error {
	var u *models.User

	if err := c.BodyParser(&u); err != nil {
		return serverError(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := u.ValidateLoginRequest(); err != nil {
		return validationError(c, err)
	}

	ctx := c.Context()

	credentials, err := h.userService.GetIDAndPassword(ctx, u.Username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return clientError(c, fiber.StatusUnauthorized, errInvalidCredentials)
		}

		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := util.CompareHashAndPassword(credentials.Password, u.Password); err != nil {
		return clientError(c, fiber.StatusUnauthorized, errInvalidCredentials)
	}

	authUser, err := h.userService.FindByID(ctx, credentials.ID)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return successResponse(c, fiber.StatusCreated, fiber.Map{
		"user": models.User{
			ID:        authUser.ID,
			Username:  authUser.Username,
			CreatedAt: authUser.CreatedAt,
			UpdatedAt: authUser.UpdatedAt,
		},
	})
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
