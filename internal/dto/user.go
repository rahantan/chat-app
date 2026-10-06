package dto

type UserResponse struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}
