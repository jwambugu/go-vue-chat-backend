package handlers

import (
	"chatapp/pkg/accesstoken"
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
		PasetoKey   string
	}

	// authHandler handles user auth
	authHandler struct {
		userService user.Service
		pasetoKey   string
	}

	authUser struct {
		ID        uint64    `json:"id,omitempty"`
		Username  string    `json:"username,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" `
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	// authUserResponse has fields returned when a user authenticates successfully
	authUserResponse struct {
		User  authUser `json:"user,omitempty"`
		Token string   `json:"token,omitempty"`
	}
)

// generateAccessToken attempts to create an access token to authenticate the user
func (h *authHandler) generateAccessToken(user *models.User) (string, error) {
	maker, err := accesstoken.NewPasetoMaker(h.pasetoKey)
	if err != nil {
		return "", err
	}

	token, err := maker.CreateToken(user, 30*time.Minute)
	if err != nil {
		return "", err
	}

	return token, nil
}

// successAuthResponse builds the response data used when a user registers or logs in
func successAuthResponse(user *models.User, token string) authUserResponse {
	return authUserResponse{
		User: authUser{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}
}

// Register adds and returns the new user created
func (h *authHandler) Register(c *fiber.Ctx) error {
	var u *models.User

	if err := c.BodyParser(&u); err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
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

	token, err := h.generateAccessToken(newUser)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return successResponse(c, fiber.StatusCreated, successAuthResponse(newUser, token))
}

// Login attempts to log in a user using the provided credentials
func (h *authHandler) Login(c *fiber.Ctx) error {
	var u *models.User

	if err := c.BodyParser(&u); err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
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

	token, err := h.generateAccessToken(authUser)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return successResponse(c, fiber.StatusOK, successAuthResponse(authUser, token))
}

// getAuthUserTokenPayload parses the value stored from the auth middleware to token payload
func getAuthUserTokenPayload(c *fiber.Ctx) *accesstoken.Payload {
	return c.Locals(accesstoken.AuthUserToken).(*accesstoken.Payload)
}

// getAuthUser returns the current auth user from the token payload
func getAuthUser(c *fiber.Ctx) *models.User {
	return getAuthUserTokenPayload(c).User
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
		pasetoKey:   opts.PasetoKey,
	}
}
