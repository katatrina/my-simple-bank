package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates a new token for a specific username and duration.
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not.
	VerifyToken(token string) (*jwt.RegisteredClaims, error)
}
