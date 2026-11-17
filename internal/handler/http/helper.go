package http

import (
	"errors"
	"realtime-chat/internal/dto/response"

	"github.com/gofiber/fiber/v2"
)

func getUserData(ctx *fiber.Ctx) (*response.User, error) {
	user, ok := ctx.Locals("user").(*response.User)
	if !ok {
		return nil, errors.New("forbidden")
	}

	return user, nil
}

func successResponse(ctx *fiber.Ctx, statusCode int, msg string, data any) error {
	return ctx.Status(statusCode).JSON(response.StandardResponse{
		Status:  true,
		Message: msg,
		Data:    data,
	})
}
