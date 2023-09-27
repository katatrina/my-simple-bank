package token

import "time"

// Maker is an interface for managing tokens.
type Maker interface {
	// CreateToken creates and sign a new token for a specific username and valid duration.
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token provided is valid or not.
	VerifyToken(token string) (*Payload, error)
}
