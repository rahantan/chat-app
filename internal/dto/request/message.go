package request

type AddMessage struct {
	// RoomID  int64  `json:"room_id" validate:"required,numeric"`
	MsgType string `json:"message_type" validate:"required"`
	Content string `json:"content" validate:"required"`
}
