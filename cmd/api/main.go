package main

import (
	"fmt"
	"log"
	"realtime-chat/config"
	httpHandler "realtime-chat/internal/handler/http"
	"realtime-chat/internal/middleware"
	mysqlRepo "realtime-chat/internal/repository/mysql"
	"realtime-chat/internal/router"
	"realtime-chat/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {

	conf := config.LoadConfig()
	conn := conf.GetConnectionDB()
	validate := conf.GetValidator()

	baseRepo := mysqlRepo.NewBaseRepository(conn)
	txRepo := mysqlRepo.NewTransactorRepository(conn)

	refreshTokenRepo := mysqlRepo.NewRefreshTokenRepository(baseRepo)
	userRepo := mysqlRepo.NewUserRepository(baseRepo)
	contactRepo := mysqlRepo.NewContactRepository(baseRepo)
	roomRepo := mysqlRepo.NewRoomRepository(baseRepo)
	roomMemberRepo := mysqlRepo.NewRoomMemberRepository(baseRepo)
	msgRepo := mysqlRepo.NewMessageRepository(baseRepo)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, &conf.Jwt)
	roomService := service.NewRoomService(roomRepo, roomMemberRepo, userRepo)
	contactService := service.NewContactService(contactRepo, userRepo)
	userService := service.NewUserService(userRepo)
	msgService := service.NewMessageService(roomRepo, roomMemberRepo, msgRepo, txRepo)

	authHandler := httpHandler.NewAuthHandler(authService, validate)
	userHandler := httpHandler.NewUserHandler(userService)
	contactHandler := httpHandler.NewContactHandler(contactService)
	roomHandler := httpHandler.NewRoomHandler(roomService)
	msgHandler := httpHandler.NewMessageHandler(msgService)

	authMiddleWare := middleware.NewAuthMiddleWare(authService, &conf.Jwt)

	dependency := router.Dependency{
		AuthHandler:    authHandler,
		AuthMiddleWare: authMiddleWare,
		ContactHandler: contactHandler,
		UserHandler:    userHandler,
		RoomHandler:    roomHandler,
		MessageHandler: msgHandler,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: httpHandler.ErrorHandler,
	})

	router.Router(app, &dependency)

	if err := app.Listen(fmt.Sprintf("%s:%s",
		conf.Server.Host,
		conf.Server.Port,
	)); err != nil {
		log.Fatal("ERROR LISTEN: ", err.Error())
	}
}
