package handlers

import (
	"chatapp/pkg/models"
	"chatapp/services/chatroom"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var (
	errInvalidCharRoomID = "Invalid chatroom id provided."
)

type (
	// ChatRoomHandlerOptions represents the options required to set up the chat room handler
	ChatRoomHandlerOptions struct {
		ChatRoomService chatroom.Service
	}

	// chatRoomHandler handles chat room interactions
	chatRoomHandler struct {
		chatRoomService chatroom.Service
	}
)

// Store creates a new models.ChatRoom
func (h *chatRoomHandler) Store(c *fiber.Ctx) error {
	var chatRoom *models.ChatRoom

	if err := c.BodyParser(&chatRoom); err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := chatRoom.ValidateStoreRequest(); err != nil {
		return validationError(c, err)
	}

	now := time.Now()
	chatRoom.CreatedAt = now
	chatRoom.UpdatedAt = now

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	chatRoom.UUID = randomUUID

	newChatRoom, err := h.chatRoomService.Create(c.Context(), chatRoom)
	if err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return successResponse(c, fiber.StatusCreated, fiber.Map{
		"chat_room": newChatRoom,
	})
}

// findChatRoomError returns the errors that occur fetching a chat room
func findChatRoomError(c *fiber.Ctx, err error) error {
	if errors.Is(err, models.ErrNoRecord) {
		return clientError(c, fiber.StatusNotFound, "Chat room not found.")
	}

	return serverError(c, fiber.StatusInternalServerError, err.Error())
}

// GetByID finds a models.ChatRoom by ID
func (h *chatRoomHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Get("id"))

	if err != nil {
		return clientError(c, fiber.StatusBadRequest, errInvalidCharRoomID)
	}

	chatRoom, err := h.chatRoomService.FindByID(c.Context(), uint64(id))
	if err != nil {
		return findChatRoomError(c, err)
	}

	return successResponse(c, fiber.StatusOK, fiber.Map{
		"chat_room": chatRoom,
	})
}

// GetByUUID finds a models.ChatRoom by UUID
func (h *chatRoomHandler) GetByUUID(c *fiber.Ctx) error {
	chatRoomUUID := c.Get("uuid")

	chatRoom, err := h.chatRoomService.FindByUUID(c.Context(), chatRoomUUID)
	if err != nil {
		return findChatRoomError(c, err)
	}

	return successResponse(c, fiber.StatusOK, fiber.Map{
		"chat_room": chatRoom,
	})
}

// Destroy soft deletes a models.ChatRoom
func (h *chatRoomHandler) Destroy(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Get("id"))

	if err != nil {
		return clientError(c, fiber.StatusBadRequest, errInvalidCharRoomID)
	}

	if err := h.chatRoomService.SoftDelete(c.Context(), uint64(id)); err != nil {
		return serverError(c, fiber.StatusInternalServerError, err.Error())
	}

	return successResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Chatroom deleted successfully.",
	})
}

// ChatRoomHandler is an interface for rooms interactions
type ChatRoomHandler interface {
	Store(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByUUID(c *fiber.Ctx) error
	Destroy(c *fiber.Ctx) error
}

// NewChatRoomHandler creates a new ChatRoomHandler
func NewChatRoomHandler(opts ChatRoomHandlerOptions) ChatRoomHandler {
	return &chatRoomHandler{
		chatRoomService: opts.ChatRoomService,
	}
}
