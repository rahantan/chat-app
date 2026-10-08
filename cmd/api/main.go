package main

import (
	"fmt"
	"log"
	"simple-auth/config"
	httpHandler "simple-auth/internal/handler/http"
	"simple-auth/internal/middleware"
	mysqlRepo "simple-auth/internal/repository/mysql"
	"simple-auth/internal/router"
	"simple-auth/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {

	conf := config.LoadConfig()
	dbConn := conf.GetConnectionDB()

	userRepo := mysqlRepo.NewUserRepository(dbConn)
	refreshTokenRepo := mysqlRepo.NewRefreshTokenRepository(dbConn)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, &conf.Jwt)
	authHandler := httpHandler.NewAuthHandler(authService)

	userHandler := httpHandler.NewUserHandlerImpl()

	middleWare := middleware.NewMiddleWare(authService, &conf.Jwt)

	route := router.Route{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		Jwt:         conf.Jwt,
		MiddleWare:  middleWare,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: httpHandler.ErrorHandler,
	})

	router.Router(app, &route)

	if err := app.Listen(fmt.Sprintf("%s:%s",
		conf.Server.Host,
		conf.Server.Port,
	)); err != nil {
		log.Fatal("ERROR LISTEN: ", err.Error())
	}
}
