package request

type CreateGroup struct {
	Name  string `json:"name"`
	Image *ImageUpload
}

type AddMemberToGroup struct {
	RoomID        int64
	MemberUserIDs []int64
}

// type CreatePrivateRoom struct {
// 	Title         string `json:"title"`
// 	ContactUserID int64  `json:"contact_user_id"`
// }

// type CreateRoom struct {
// 	Type string
// }
