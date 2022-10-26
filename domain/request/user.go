package request

type GetUser struct {
	ID    uint64 `json:"id,omitempty" validate:"omitempty,numeric"`
	Email string `json:"email,omitempty" validate:"omitempty,email,min=6,max=200"`
}

type User struct {
	Name     string `json:"name" validate:"required,min=3,max=200"`
	Email    string `json:"email" validate:"required,email,min=6,max=200"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}
