package middleware

import (
	"errors"
	"simple-auth/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("access_token")
		if tokenStr == "" {
			return errors.New("token is required")
		}

		userResponse, err := pkg.ValidationJWTToken(secretKey, tokenStr)
		if err != nil {
			return err
		}

		c.Locals("user", *userResponse)

		return c.Next()
	}
}
