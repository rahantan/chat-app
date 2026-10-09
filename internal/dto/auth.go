package dto

type RegisterRequest struct {
	UserName    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Pass        string `json:"password" validate:"required,min=5"`
	ConfirmPass string `json:"confirm_password" validate:"required,eqfield=Pass"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"password" validate:"required"`
}

type AuthResponse struct {
	RefreshToken string
	AccessToken  string
	UserResponse
}
