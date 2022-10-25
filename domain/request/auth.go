package request

type Login struct {
	Email    string `json:"email" validate:"required,email,min=6,max=200"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

type Logout struct{}
