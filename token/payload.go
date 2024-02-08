package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// Payload contains the payload data of the token.
type Payload struct {
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload for a specific username and duration.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		jwt.RegisteredClaims{
			Issuer:    "simplebank",
			Subject:   username,
			Audience:  jwt.ClaimStrings{"client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID.String(),
		},
	}

	return payload, nil
}
