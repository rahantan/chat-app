package middleware

import (
	"realtime-chat/config"
	"realtime-chat/internal/pkg"
	"realtime-chat/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleWare interface {
	Authentication(c *fiber.Ctx) error
}

type authMiddleWareImpl struct {
	service.AuthService
	*config.Jwt
}

func NewAuthMiddleWare(auth service.AuthService, jwtConf *config.Jwt) AuthMiddleWare {
	return &authMiddleWareImpl{
		AuthService: auth,
		Jwt:         jwtConf,
	}
}
func (m *authMiddleWareImpl) Authentication(c *fiber.Ctx) error {
	tokenStr := c.Cookies("access_token")
	if tokenStr == "" {
		return pkg.ErrTokenRequired
	}

	user, err := m.AuthService.VerifyToken(m.Jwt.ACCESS_TOKEN_SECRET_KEY, tokenStr)
	if err != nil {
		return err
	}

	c.Locals("user", user)

	return c.Next()
}
