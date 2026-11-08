package response

type ChatGroup struct {
	RoomID    int64         `json:""`
	Type      string        `json:""`
	Name      string        `json:""`
	ImagePath string        `json:""`
	Members   []GroupMember `json:""`
	CreatedAt string        `json:""`
	UpdatedAt string        `json:""`
}

type GroupMember struct {
	MemberID int64  `json:"member_id"`
	UserID   int64  `json:"user_id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}
