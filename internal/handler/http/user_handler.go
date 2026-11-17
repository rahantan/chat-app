package http

import (
	"context"
	"os"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Get(ctx *fiber.Ctx) error
	GetAvatar(ctx *fiber.Ctx) error
}

type userHandlerImpl struct {
	service.UserService
}

func NewUserHandler(user service.UserService) UserHandler {
	return &userHandlerImpl{
		UserService: user,
	}
}

func (userHandler *userHandlerImpl) Get(ctx *fiber.Ctx) error {

	userID, err := ctx.ParamsInt("userID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	result, err := userHandler.UserService.GetById(c, int64(userID))
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get user", result)
}

func (userHandler *userHandlerImpl) GetAvatar(ctx *fiber.Ctx) error {
	fileName := ctx.Query("src")
	if fileName == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	filePath := pkg.UserImagePath + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.SendFile(pkg.UserImagePath + fileName)
}
