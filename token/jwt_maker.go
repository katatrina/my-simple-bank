package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

// NewJWTMaker returns a new JWTMaker.
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	// Create a new jwtToken
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign the jwtToken
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// Parse the token
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		// TODO: Handle error more gracefully
		return nil, err
	}

	// Check if the token is valid
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return payload, nil
}
