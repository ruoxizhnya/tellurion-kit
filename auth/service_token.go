package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IssueServiceToken mints a long-lived access token for a service account.
//
// Unlike IssueTokens it is not tied to a user row: subject is caller-chosen
// (for example "service:risk"), role is an opaque caller-chosen string (for
// example "service"), and ttl is explicit. Only an access token is produced —
// there is no refresh half, so this token never goes through Refresh.
//
// Rotation and distribution belong to the caller (core QT-71): mint the token
// from an operator-controlled process, inject it through the deployment's
// secret mechanism, and rotate it on the caller's schedule. This package never
// stores tokens, keeps no revocation list, and knows no service registry.
//
// Empty (or whitespace-only) subject, empty (or whitespace-only) role, and
// non-positive ttl are rejected with ErrInvalidServiceToken. A disabled
// service (empty JWT secret) returns ErrTokensDisabled, like IssueTokens.
func (s *TokenService) IssueServiceToken(subject, role string, ttl time.Duration) (string, error) {
	if !s.Enabled() {
		return "", ErrTokensDisabled
	}
	if strings.TrimSpace(subject) == "" {
		return "", fmt.Errorf("%w: subject must not be empty", ErrInvalidServiceToken)
	}
	if strings.TrimSpace(role) == "" {
		return "", fmt.Errorf("%w: role must not be empty", ErrInvalidServiceToken)
	}
	if ttl <= 0 {
		return "", fmt.Errorf("%w: ttl must be positive (got %s)", ErrInvalidServiceToken, ttl)
	}

	now := time.Now().UTC()
	claims := Claims{
		Username: subject,
		Role:     role,
		Kind:     KindAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return s.mint(claims)
}
