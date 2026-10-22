package auth

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// Token kinds. Kind is a discriminator so a refresh token can never be
// replayed as an access token.
const (
	KindAccess  = "access"
	KindRefresh = "refresh"
)

// Claims is the JWT payload this package signs. It extends RegisteredClaims
// with the minimum identity downstream handlers need: username, role, and the
// token kind.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Kind     string `json:"kind"` // KindAccess | KindRefresh
	jwt.RegisteredClaims
}

// UserIDInt64 returns the subject as int64. Returns an error if the subject
// is empty or not a valid integer.
func (c *Claims) UserIDInt64() (int64, error) {
	if c.Subject == "" {
		return 0, fmt.Errorf("auth: token has empty subject")
	}
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("auth: subject is not an int64: %w", err)
	}
	return id, nil
}

func formatUserID(id int64) string {
	return strconv.FormatInt(id, 10)
}
