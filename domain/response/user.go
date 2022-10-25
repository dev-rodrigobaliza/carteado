package response

type User struct {
	ID      uint64 `json:"id,omitempty"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}
