package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker.
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d  character", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot create token ID: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "simple_bank",
		Subject:   username,
		Audience:  jwt.ClaimStrings{"client"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        tokenID.String(),
	})

	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*jwt.RegisteredClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}

	return claims, nil
}
