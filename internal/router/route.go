package router

import (
	"simple-auth/config"
	httpHandler "simple-auth/internal/handler/http"
	"simple-auth/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	httpHandler.AuthHandler
	httpHandler.UserHandler
	middleware.MiddleWare
	config.Jwt
}

func Router(app *fiber.App, route *Route) {

	api := app.Group("/api")

	public := api.Group("/public")
	public.Post("/login", route.AuthHandler.Login)
	public.Post("/register", route.AuthHandler.Register)
	public.Post("/refresh-token", route.AuthHandler.RefreshToken)

	// SECURE
	protect := api.Group("/protect", route.MiddleWare.AuthMiddleware)

	protect.Post("/log-out", route.AuthHandler.LogOut)
	protect.Get("/me", route.UserHandler.GetUser)
}
