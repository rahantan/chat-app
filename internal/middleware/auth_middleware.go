package middleware

import (
	"realtime-chat/config"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"

	"github.com/gofiber/fiber/v2"
)

type MiddleWare interface {
	AuthMiddleware(c *fiber.Ctx) error
}

type MiddleWareImpl struct {
	service.AuthService
	*config.Jwt
}

func NewMiddleWare(auth service.AuthService, jwtConf *config.Jwt) MiddleWare {
	return &MiddleWareImpl{
		AuthService: auth,
		Jwt:         jwtConf,
	}
}
func (m *MiddleWareImpl) AuthMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Cookies("access_token")
	if tokenStr == "" {
		return pkg.ErrTokenRequired
	}

	userResponse, err := m.AuthService.VerifyToken(m.Jwt.ACCESS_TOKEN_SECRET_KEY, tokenStr)
	if err != nil {
		return err
	}

	c.Locals("user", *userResponse)

	return c.Next()
}
