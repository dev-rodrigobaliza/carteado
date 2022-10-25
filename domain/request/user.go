package request

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=200"`
	Email    string `json:"email" validate:"required,email,min=6,max=200"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}
