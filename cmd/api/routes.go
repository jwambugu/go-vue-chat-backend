package main

import (
	"chatapp/cmd/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func (app *application) routes() *fiber.App {
	fiberApp := fiber.New()

	registerFiberMiddleware(fiberApp)

	v1 := fiberApp.Group("api/v1")

	auth := v1.Group("/auth")
	authHandler := handlers.NewAuthHandler(handlers.AuthHandlerOptions{
		UserService: app.userService,
		PasetoKey:   app.config.PasetoKey,
	})

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	chatRooms := v1.Group("/chat-rooms").Use(app.authMiddleware())
	chatRoomsHandler := handlers.NewChatRoomHandler(handlers.ChatRoomHandlerOptions{
		ChatRoomService: app.chatroomService,
	})

	chatRooms.Get("/", chatRoomsHandler.Index)
	chatRooms.Post("/", chatRoomsHandler.Store)
	chatRooms.Get("/:id", chatRoomsHandler.Show)
	chatRooms.Get("/:uuid/uuid", chatRoomsHandler.GetByUUID)
	chatRooms.Delete("/:id", chatRoomsHandler.Destroy)

	return fiberApp
}
