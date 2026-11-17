package request

type AddContact struct {
	ContactUserID int64 `json:"contact_user_id" validate:"required,numeric"`
}
