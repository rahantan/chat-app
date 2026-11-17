package http

import (
	"context"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler interface {
	List(ctx *fiber.Ctx) error
	Send(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type messageHandlerImpl struct {
	service.MessageService
}

func NewMessageHandler(msg service.MessageService) MessageHandler {
	return &messageHandlerImpl{
		MessageService: msg,
	}
}

func (msgHandler *messageHandlerImpl) List(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	roomID, err := ctx.ParamsInt("roomID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	result, err := msgHandler.MessageService.ListMessage(c, user.ID, int64(roomID))
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get messages", result)
}

func (msgHandler *messageHandlerImpl) Send(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	roomID, err := ctx.ParamsInt("roomID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	var reqBody request.AddMessage
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := msgHandler.MessageService.SendMessage(c, user.ID, int64(roomID), &reqBody); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (msgHandler *messageHandlerImpl) Delete(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}
	roomID, err := ctx.ParamsInt("roomID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	msgID, err := ctx.ParamsInt("messageID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := msgHandler.MessageService.DeleteMessage(c, user.ID, int64(roomID), int64(msgID)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
