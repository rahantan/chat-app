package response

type Group struct {
}

type RoomLastMsg struct {
	SenderID    int64
	Content     string
	TypeContent string
	CreatedAt   string
}
type Room struct {
	RoomID    int64
	Type      string
	Name      string
	ImagePath string
	RoomLastMsg
	UpdatedAt string
}
