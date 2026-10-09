package http

import (
	"context"
	"errors"
	"log"
	"simple-auth/internal/dto"
	"simple-auth/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errApp, ok := err.(*pkg.AppErr); ok {
		log.Println(errApp.Message)
		log.Println(errApp.Details)

		return c.Status(pkg.GetStatusCode(errApp)).JSON(dto.StandardResponse{
			Status:  false,
			Message: errApp.Message,
			Error:   errApp.Details,
		})
	}

	log.Println(err.Error())

	if errors.Is(context.DeadlineExceeded, err) {
		return c.Status(fiber.StatusRequestTimeout).JSON(dto.StandardResponse{
			Status:  false,
			Message: "request time out",
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(dto.StandardResponse{
		Status:  false,
		Message: "internal server error",
	})
}
