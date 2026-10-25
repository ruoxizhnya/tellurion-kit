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

// ServiceRole is the name of the service-to-service principal (裁-17).
//
// It is NOT a role enum: this package deliberately does not define the role
// space — callers pass role strings through IssueServiceToken (KT-09 ③-④).
// It is named here for one reason only: a service token's subject is
// caller-chosen and non-numeric ("service:risk"), so the int64 subject check
// in the gin middleware and in auth/httpadapter must exempt it. Without that
// exemption every legitimate service token 401s at the API boundary even
// though its signature and role are valid (verified 2026-10-25 by minting one
// with IssueServiceToken and feeding it through httpadapter.Authenticate).
const ServiceRole = "service"

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
