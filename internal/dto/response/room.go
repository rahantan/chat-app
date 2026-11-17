package response

type Group struct {
}

type RoomLastMsg struct {
	UserID      int64  `json:"user_id"`
	UserName    string `json:"username"`
	Content     string `json:"content"`
	TypeContent string `json:"type_content"`
	CreatedAt   string `json:"created_at"`
}

type Member struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type Room struct {
	RoomID    int64        `json:"room_id"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	ImageUrl  string       `json:"image_url"`
	LastMsg   *RoomLastMsg `json:"last_message"`
	Members   []Member     `json:"members,omitempty"`
	UpdatedAt string       `json:"updated_at"`
}
