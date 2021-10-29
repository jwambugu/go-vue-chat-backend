package main

import (
	"chatapp/cmd/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func (app *application) routes() *fiber.App {
	fiberApp := fiber.New()

	v1 := fiberApp.Group("api/v1")

	auth := v1.Group("/auth")
	authHandler := handlers.NewAuthHandler(handlers.AuthHandlerOptions{
		UserService: app.userService,
	})

	auth.Post("/register", authHandler.Register)
	return fiberApp
}
