package config

type Security struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpireTime  string `json:"expire_time,omitempty"`
}
