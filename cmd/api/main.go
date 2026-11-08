package main

import (
	"fmt"
	"log"
	"reflect"
	"realtime-chat/config"
	httpHandler "realtime-chat/internal/handler/http"
	"realtime-chat/internal/middleware"
	mysqlRepo "realtime-chat/internal/repository/mysql"
	"realtime-chat/internal/router"
	"realtime-chat/internal/service"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {

	conf := config.LoadConfig()
	dbConn := conf.GetConnectionDB()
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Tag.Get("json")
		if tag == "-" {
			return ""
		}
		return strings.Split(tag, ",")[0]
	})

	refreshTokenRepo := mysqlRepo.NewRefreshTokenRepository(dbConn)
	userRepo := mysqlRepo.NewUserRepository(dbConn)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, &conf.Jwt)
	authHandler := httpHandler.NewAuthHandler(authService, validate)

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
