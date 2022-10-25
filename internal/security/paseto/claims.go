package paseto

import "github.com/vk-rv/pvx"

type TWSClaims struct {
	pvx.RegisteredClaims
}

// NewTWSClaims creates a new TWSClaims
func NewTWSClaims(registeredClaims pvx.RegisteredClaims) *TWSClaims {
	return &TWSClaims{
		RegisteredClaims: registeredClaims,
	}
}