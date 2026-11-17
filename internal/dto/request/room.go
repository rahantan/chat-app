package request

type CreateGroup struct {
	Name  string `json:"name" form:"name" validate:"required"`
	Image *ImageUpload
}

type AddMemberToGroup struct {
	UserIDs []int64 `json:"user_ids"`
}

type GetOrCreatePrivateRoom struct {
	TargetUserID int64 `json:"target_user_id"`
}

// type CreatePrivateRoom struct {
// 	Title         string `json:"title"`
// 	ContactUserID int64  `json:"contact_user_id"`
// }

// type CreateRoom struct {
// 	Type string
// }
