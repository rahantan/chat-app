package request

type AddMessage struct {
	RoomID  int64
	MsgType string
	Content string
}
