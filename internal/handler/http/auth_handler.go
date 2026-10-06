package http

import (
	"context"
	"errors"
	"simple-auth/config"
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
	*config.Jwt
}

func NewAuthHandler(authService service.AuthService, jwtConf *config.Jwt) AuthHandler {
	return &authHandlerImpl{
		AuthService: authService,
		Jwt:         jwtConf,
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
		Path:     "/api/public/refresh-token",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
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

	user, err := authHandler.AuthService.Login(ctx, &reqBody)
	if err != nil {
		return err
	}

	refreshToken, err := pkg.CreateJWTToken(authHandler.Jwt.REFRESH_TOKEN_SECRET_KEY, pkg.RefreshTokenExpireTime, user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/public/refresh-token",
		Secure:   true,
		HTTPOnly: true,
		Expires:  time.Now().Add(pkg.RefreshTokenExpireTime),
	})

	accessToken, err := pkg.CreateJWTToken(authHandler.Jwt.ACCESS_TOKEN_SECRET_KEY, pkg.AccessTokenExpireTime, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.StandardResponse{
		Status:  true,
		Message: "SUCCESS LOGIN",
		Data: map[string]any{
			"access_token": accessToken,
			"user":         user,
		}},
	)
}

func (authHandler *authHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	tokenStr := c.Cookies("refresh_token")

	if tokenStr == "" {
		return errors.New("token is required")
	}

	user, err := pkg.ValidationJWTToken(authHandler.Jwt.REFRESH_TOKEN_SECRET_KEY, tokenStr)
	if err != nil {
		return err
	}

	token, err := pkg.CreateJWTToken(authHandler.Jwt.ACCESS_TOKEN_SECRET_KEY, pkg.AccessTokenExpireTime, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.StandardResponse{
		Status:  true,
		Message: "SUCCESS REFRESH TOKEN",
		Data: map[string]string{
			"access_token": token,
		},
	})
}
