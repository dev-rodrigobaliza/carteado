package config

type HTTP struct {
	Address      string       `json:"address,omitempty"`
	CORS         string       `json:"cors,omitempty"`
	IdleTimeout  string       `json:"idle_timeout,omitempty"`
	ReadTimeout  string       `json:"read_timeout,omitempty"`
	WriteTimeout string       `json:"write_timeout,omitempty"`
	Limiter      *HTTPLimiter `json:"limiter,omitempty"`
}

type HTTPLimiter struct {
	Enabled     bool     `json:"enabled,omitempty"`
	MaxRequests int      `json:"max_requests,omitempty"`
	Expiration  string   `json:"expiration,omitempty"`
	AllowedIPs  []string `json:"allowed_ips,omitempty"`
}
