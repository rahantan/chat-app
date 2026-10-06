package main

import (
	"fmt"
	"log"
	"simple-auth/config"
	httpHandler "simple-auth/internal/handler/http"
	mysqlRepo "simple-auth/internal/repository/mysql"
	"simple-auth/internal/router"
	"simple-auth/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {

	conf := config.LoadConfig()
	dbConn := conf.GetConnectionDB()

	userRepo := mysqlRepo.NewUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)
	authHandler := httpHandler.NewAuthHandler(authService, &conf.Jwt)

	userHandler := httpHandler.NewUserHandlerImpl()

	route := router.Route{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		Jwt:         conf.Jwt,
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
