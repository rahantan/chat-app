package response

type Contact struct {
	ID            int64  `json:"id"`
	ContactUserID int64  `json:"contact_user_id"`
	UserName      string `json:"username"`
	ImagePath     string `json:"image_path"`
	// RoomID        int64  `json:"room_id"`
}
