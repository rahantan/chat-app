package request

type Register struct {
	UserName    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Pass        string `json:"password" validate:"required,min=5"`
	ConfirmPass string `json:"confirm_password" validate:"required,eqfield=Pass"`
}

type Login struct {
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"password" validate:"required"`
}
