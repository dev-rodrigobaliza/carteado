package config

type dbConnectionPool struct {
	Idle     int `json:"idle,omitempty"`
	Open     int `json:"open,omitempty"`
	Lifetime int `json:"lifetime,omitempty"`
}

type Database struct {
	Type     string            `json:"type,omitempty"`
	Host     string            `json:"host,omitempty"`
	Port     string            `json:"port,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Name     string            `json:"name,omitempty"`
	SSL      bool              `json:"ssl,omitempty"`
	Pool     *dbConnectionPool `json:"pool,omitempty"`
}
