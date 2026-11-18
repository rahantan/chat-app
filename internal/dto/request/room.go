package request

type CreateGroup struct {
	Name  string `json:"name" form:"name" validate:"required"`
	Image *ImageUpload
}

type RemoveMembersFromGroup struct {
	UserIDs []int64 `json:"user_ids"`
}

type AddMembersToGroup struct {
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
