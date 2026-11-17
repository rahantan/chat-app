package http

import (
	"context"
	"errors"
	"log"

	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err, ok := err.(*response.Error); ok {
		log.Println(err.Msg)
		log.Println(err.Details)

		return c.Status(pkg.GetStatusCode(err)).JSON(response.StandardResponse{
			Status:    false,
			Message:   err.Msg,
			Code:      err.Code,
			ErrDetail: err.Details,
		})
	}

	log.Println(err.Error())
	if errors.Is(context.DeadlineExceeded, err) {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.StandardResponse{
			Status:  false,
			Code:    pkg.ErrCodeTimeOut,
			Message: "request time out",
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response.StandardResponse{
		Status:  false,
		Message: "internal server error",
		Code:    pkg.ErrCodeInternalError,
	})
}
