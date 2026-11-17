package router

import (
	handler "realtime-chat/internal/handler/http"
	"realtime-chat/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Dependency struct {
	handler.AuthHandler
	handler.ContactHandler
	handler.UserHandler
	handler.MessageHandler
	handler.RoomHandler
	middleware.AuthMiddleWare
}

func Router(app *fiber.App, depend *Dependency) {

	api := app.Group("/api")

	public := api.Group("/public")
	depend.public(public)

	protect := api.Group("/protect", depend.AuthMiddleWare.Authentication)
	protect.Post("/log-out", depend.AuthHandler.LogOut)

	me := protect.Group("/me")
	depend.contact(me)
	depend.room(me)
	// depend.message(me)

}

func (depend Dependency) public(route fiber.Router) {
	authHandler := depend.AuthHandler
	route.Post("/login", authHandler.Login)                // done
	route.Post("/register", authHandler.Register)          // done
	route.Post("/refresh-token", authHandler.RefreshToken) // done

	user := depend.UserHandler
	route.Get("/users/avatar", user.GetAvatar)
	route.Get("/users/:userID", user.Get)
}

func (depend Dependency) contact(route fiber.Router) {
	contact := depend.ContactHandler

	route.Delete("/contacts/:contactID", contact.Delete)
	route.Get("/contacts", contact.List)
	route.Post("/contacts", contact.Add)
}

// func (depend Dependency) message(route fiber.Router) {
// 	msg := depend.MessageHandler

// 	route.Delete("/messages/:messageID", msg.Delete)
// 	route.Post("/messages", msg.Send)
// 	route.Get("/messages", msg.List)
// }

func (depend Dependency) room(route fiber.Router) {
	room := depend.RoomHandler
	msg := depend.MessageHandler

	route.Delete("/rooms/:roomID/messages/:messageID", msg.Delete)
	route.Post("/rooms/:roomID/messages", msg.Send)
	route.Get("/rooms/:roomID/messages", msg.List)

	route.Delete("/rooms/:roomID/group", room.DeleteGroup)
	route.Post("/rooms/:roomID/group", room.AddMembersToGroup)
	route.Post("rooms/group", room.CreateGroup)

	route.Post("/rooms/private", room.GetOrCreatePrivateRoom)

	route.Get("/rooms/:roomID", room.GetRoomDetail)
	route.Get("/rooms", room.ListMyRooms)
}
