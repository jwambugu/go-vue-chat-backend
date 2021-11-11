package main

import (
	"chatapp/pkg/accesstoken"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"strings"
)

func registerFiberMiddleware(app *fiber.App) {
	app.Use(compress.New(),
		cors.New(),
		requestid.New(),
		logger.New(),
		recover.New(),
	)
}

// authMiddleware attempts to verify the access token provided before completing the request
func (app *application) authMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		requestAccessToken := strings.Split(c.Get("Authorization"), " ")

		if len(requestAccessToken) == 0 || requestAccessToken[0] != "Bearer" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bearer authorization header is required.",
			})
		}

		maker, err := accesstoken.NewPasetoMaker(app.config.PasetoKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tokenPayload, err := maker.VerifyToken(requestAccessToken[1])
		if err != nil {
			if errors.Is(err, accesstoken.ErrInvalidToken) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": accesstoken.ErrInvalidToken.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals(accesstoken.AuthUserToken, tokenPayload)
		return c.Next()
	}
}
