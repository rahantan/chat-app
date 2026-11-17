package http

import (
	"context"

	"realtime-chat/internal/dto/request"
	"realtime-chat/internal/dto/response"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	LogOut(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
}
type authHandlerImpl struct {
	service.AuthService
	*validator.Validate
}

func NewAuthHandler(authService service.AuthService, v *validator.Validate) AuthHandler {
	return &authHandlerImpl{
		AuthService: authService,
		Validate:    v,
	}
}

func (authHandler *authHandlerImpl) Register(ctx *fiber.Ctx) error {
	var reqBody request.Register
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ErrInvalidPayload
	}

	if err := authHandler.Validate.Struct(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	err := authHandler.AuthService.Register(c, &reqBody)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (authHandler *authHandlerImpl) Login(ctx *fiber.Ctx) error {
	var reqBody request.Login
	if err := ctx.BodyParser(&reqBody); err != nil {
		return pkg.ErrInvalidPayload
	}

	if err := authHandler.Validate.Struct(&reqBody); err != nil {
		return pkg.ValidationError(err)
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	result, err := authHandler.AuthService.Login(c, &reqBody)
	if err != nil {
		return err
	}

	authHandler.setAuthCookies(ctx, result)

	return successResponse(ctx, fiber.StatusOK, "success login", result.User)
}

func (authHandler *authHandlerImpl) LogOut(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     pkg.RefreshTokenCookiePath,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
		Secure:   false,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     pkg.AccessTokenCookiePath,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: false,
		Secure:   false,
	})
	return ctx.SendStatus(fiber.StatusOK)
}

func (authHandler *authHandlerImpl) RefreshToken(ctx *fiber.Ctx) error {
	tokenStr := ctx.Cookies("refresh_token")
	if tokenStr == "" {
		return pkg.ErrTokenRequired
	}

	c, cancel := context.WithTimeout(ctx.UserContext(), 5*time.Second)
	defer cancel()

	result, err := authHandler.AuthService.RefreshToken(c, tokenStr)
	if err != nil {
		return err
	}

	authHandler.setAuthCookies(ctx, result)

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (authHandler *authHandlerImpl) setAuthCookies(ctx *fiber.Ctx, auth *response.Auth) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    auth.RefreshToken,
		Path:     pkg.RefreshTokenCookiePath,
		Secure:   false,
		HTTPOnly: true,
		Expires:  time.Now().Add(pkg.RefreshTokenExpireTime),
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    auth.AccessToken,
		Path:     pkg.AccessTokenCookiePath,
		Secure:   false,
		HTTPOnly: true,
		Expires:  time.Now().Add(pkg.AccessTokenExpireTime),
	})

}
