package http

import (
	"context"
	"errors"
	"simple-auth/internal/dto"
	"simple-auth/internal/pkg"
	"simple-auth/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	LogOut(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}
type authHandlerImpl struct {
	service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandlerImpl{
		AuthService: authService,
	}
}

func (authHandler *authHandlerImpl) Register(c *fiber.Ctx) error {
	var reqBody dto.RegisterRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	err := authHandler.AuthService.Register(ctx, &reqBody)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (authHandler *authHandlerImpl) LogOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     pkg.RefreshTokenCookiePath,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
		Secure:   false,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     pkg.AccessTokenCookiePath,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
		Secure:   false,
	})
	return c.SendStatus(fiber.StatusOK)
}

func (authHandler *authHandlerImpl) Login(c *fiber.Ctx) error {
	var reqBody dto.LoginRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	result, err := authHandler.AuthService.Login(ctx, &reqBody)
	if err != nil {
		return err
	}

	authHandler.setAuthCookies(c, result)

	return c.Status(fiber.StatusOK).JSON(dto.StandardResponse{
		Status:  true,
		Message: "SUCCESS LOGIN",
		Data: map[string]any{
			"user": result.UserResponse,
		}},
	)
}

func (authHandler *authHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	tokenStr := c.Cookies("refresh_token")
	if tokenStr == "" {
		return errors.New("token is required")
	}

	result, err := authHandler.AuthService.RefreshToken(c.UserContext(), tokenStr)
	if err != nil {
		return err
	}

	authHandler.setAuthCookies(c, result)

	return c.SendStatus(fiber.StatusOK)
}

func (authHandler *authHandlerImpl) setAuthCookies(c *fiber.Ctx, auth *dto.AuthResponse) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    auth.RefreshToken,
		Path:     pkg.RefreshTokenCookiePath,
		Secure:   false,
		HTTPOnly: true,
		Expires:  time.Now().Add(pkg.RefreshTokenExpireTime),
	})

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    auth.AccessToken,
		Path:     pkg.AccessTokenCookiePath,
		Secure:   false,
		HTTPOnly: true,
		Expires:  time.Now().Add(pkg.AccessTokenExpireTime),
	})

}
