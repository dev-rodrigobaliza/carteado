package paseto

import (
	"time"

	"github.com/vk-rv/pvx"
)

const (
	SYMMETRIC_KEY_SIZE = 32
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	assert       []byte
	paseto       *pvx.ProtoV4Local
	symmetricKey *pvx.SymKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey, assert string) *PasetoMaker {
	if len(symmetricKey) < SYMMETRIC_KEY_SIZE {
		for i := len(symmetricKey); i < SYMMETRIC_KEY_SIZE; i++ {
			symmetricKey += "0"
		}
	}
	symKey := pvx.NewSymmetricKey([]byte(symmetricKey), pvx.Version4)

	token := &PasetoMaker{
		assert:       []byte(assert),
		paseto:       pvx.NewPV4Local(),
		symmetricKey: symKey,
	}

	return token
}

// CreateToken creates a new token for a specific subject (uuid) and duration
func (m *PasetoMaker) CreateToken(subject string, duration time.Duration) (string, error) {
	// make data
	exp := time.Now().Add(duration)

	claims := &TWSClaims{
		RegisteredClaims: pvx.RegisteredClaims{
			Subject:    subject,
			Expiration: &exp,
		},
	}

	token, err := m.paseto.Encrypt(m.symmetricKey, claims, pvx.WithAssert([]byte(m.assert)))
	if err != nil {
		return ": ", err
	}

	return token, nil
}

// VerifyToken checks if the token is valid or not
func (m *PasetoMaker) VerifyToken(token string) (string, error) {
	claims := &TWSClaims{}
	err := m.paseto.Decrypt(token, m.symmetricKey, pvx.WithAssert([]byte(m.assert))).ScanClaims(claims)
	if err != nil {
		return ": ", err
	}

	return claims.Subject, nil
}
