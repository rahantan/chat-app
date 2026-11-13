package response

type Auth struct {
	RefreshToken string
	AccessToken  string
	User
}
