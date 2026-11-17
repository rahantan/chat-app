package response

type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
}
