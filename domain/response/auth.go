package response

import "time"

type Login struct {
	User        *User       `json:"user"`
	AccessToken AccessToken `json:"access_token"`
}

type AccessToken struct {
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"valid_until"`
}
