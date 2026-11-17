package response

type Contact struct {
	ID            int64  `json:"id"`
	ContactUserID int64  `json:"contact_user_id"`
	UserName      string `json:"username"`
	Avatar        string `json:"avatar"`
	// RoomID        int64  `json:"room_id"`
}
