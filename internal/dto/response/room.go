package response

type Group struct {
}

type RoomLastMsg struct {
	SenderName  string `json:""`
	Content     string `json:"content"`
	TypeContent string `json:"type_content"`
	CreatedAt   string `json:"created_at"`
}

type Member struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}

type Room struct {
	RoomID    int64       `json:"room_id"`
	Type      string      `json:"type"`
	Name      string      `json:"name"`
	ImagePath string      `json:"image_path"`
	LastMsg   RoomLastMsg `json:"last_message"`
	Members   []Member    `json:"members"`
	UpdatedAt string      `json:"updated_at"`
}

// type RoomResponse struct {
// 	RoomID      int64
// 	Type        string
// 	Name        string
// 	ImagePath   string
// 	LastMessage string
// }
