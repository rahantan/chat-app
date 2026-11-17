package http

import (
	"context"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ContactHandler interface {
	// GetContact(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Add(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type contactHandlerImpl struct {
	service.ContactService
}

func NewContactHandler(contact service.ContactService) ContactHandler {
	return &contactHandlerImpl{
		ContactService: contact,
	}
}

func (contactHandler *contactHandlerImpl) Delete(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	contactID, err := ctx.ParamsInt("contactID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := contactHandler.ContactService.Delete(c, user.ID, int64(contactID)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
func (contactHandler *contactHandlerImpl) List(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	results, err := contactHandler.ContactService.List(c, user.ID)
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get my contacts", results)
}

func (contactHandler *contactHandlerImpl) Add(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	var reqBody request.AddContact
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := contactHandler.ContactService.Add(c, user.ID, &reqBody); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
