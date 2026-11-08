package http

import (
	"realtime-chat/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	GetUser(c *fiber.Ctx) error
}

type userHandlerImpl struct {
}

func NewUserHandlerImpl() UserHandler {
	return &userHandlerImpl{}
}

func (userHandler *userHandlerImpl) GetUser(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(dto.UserResponse)
	return c.Status(fiber.StatusOK).JSON(dto.StandardResponse{
		Status:  true,
		Message: "SUCCESS GET USER",
		Data:    user,
	})
}
