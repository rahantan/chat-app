package http

import (
	"context"
	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler interface {
	ListMyRooms(ctx *fiber.Ctx) error
	GetRoomDetail(ctx *fiber.Ctx) error

	CreateGroup(ctx *fiber.Ctx) error
	AddMembersToGroup(ctx *fiber.Ctx) error
	DeleteGroup(ctx *fiber.Ctx) error

	GetOrCreatePrivateRoom(ctx *fiber.Ctx) error
}

type roomHandlerImpl struct {
	service.RoomService
}

func NewRoomHandler(room service.RoomService) RoomHandler {
	return &roomHandlerImpl{
		RoomService: room,
	}
}
func (roomHandler *roomHandlerImpl) ListMyRooms(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	results, err := roomHandler.RoomService.ListMyRooms(c, user.ID)
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get my rooms", results)
}
func (roomHandler *roomHandlerImpl) GetRoomDetail(ctx *fiber.Ctx) error {
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

	result, err := roomHandler.RoomService.GetRoomDetail(c, user.ID, int64(roomID))
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get room detail", result)
}

func (roomHandler *roomHandlerImpl) CreateGroup(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}
	var reqBody request.CreateGroup
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	var imageGroup *request.ImageUpload

	image, err := ctx.FormFile("image")
	if err == nil {
		fileContent, err := image.Open()
		if err != nil {
			return err
		}
		defer fileContent.Close()

		imageGroup = &request.ImageUpload{
			FileName: image.Filename,
			Size:     image.Size,
			Content:  fileContent,
		}
	} else {
		imageGroup = nil
	}

	reqBody.Image = imageGroup

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := roomHandler.RoomService.CreateGroupRoom(c, user.ID, &reqBody); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
func (roomHandler *roomHandlerImpl) AddMembersToGroup(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	roomID, err := ctx.ParamsInt("roomID")
	if err != nil {
		return pkg.ErrInvalidParamID
	}

	var reqBody request.AddMemberToGroup
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	if err := roomHandler.RoomService.AddMembersToGroup(c, user.ID, int64(roomID), &reqBody); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
func (roomHandler *roomHandlerImpl) DeleteGroup(ctx *fiber.Ctx) error {
	return nil
}

func (roomHandler *roomHandlerImpl) GetOrCreatePrivateRoom(ctx *fiber.Ctx) error {
	user, err := getUserData(ctx)
	if err != nil {
		return err
	}

	var reqBody request.GetOrCreatePrivateRoom
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	roomID, err := roomHandler.RoomService.GetOrCreatePrivateRoom(c, user.ID, &reqBody)
	if err != nil {
		return err
	}

	return successResponse(ctx, fiber.StatusOK, "success get", map[string]any{
		"room_id": roomID,
	})
}
