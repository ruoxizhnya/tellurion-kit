package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueServiceToken_RoundTrip(t *testing.T) {
	t.Parallel()
	s := testService("service-secret")

	const ttl = 30 * 24 * time.Hour
	token, err := s.IssueServiceToken("service:risk", "service", ttl)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := s.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, KindAccess, claims.Kind, "service tokens must be access tokens")
	assert.Equal(t, "service", claims.Role, "role is the caller-supplied string, not an enum")
	assert.Equal(t, "service:risk", claims.Subject)
	assert.Equal(t, "service:risk", claims.Username, "username mirrors the subject so audit consumers get a non-empty actor")
	assert.Equal(t, "tellurion", claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
	require.NotNil(t, claims.IssuedAt)
	assert.Equal(t, ttl, claims.ExpiresAt.Sub(claims.IssuedAt.Time),
		"the explicit ttl, not Config.AccessTokenTTL, must drive expiry")

	// No refresh half: the same token must be rejected by ParseRefresh.
	_, err = s.ParseRefresh(token)
	assert.ErrorIs(t, err, ErrInvalidToken, "service tokens are access-only; no refresh token is minted")
}

func TestIssueServiceToken_Expires(t *testing.T) {
	t.Parallel()
	s := testService("service-secret")

	token, err := s.IssueServiceToken("service:data", "service", time.Millisecond)
	require.NoError(t, err)

	time.Sleep(20 * time.Millisecond)

	_, err = s.Parse(token)
	assert.ErrorIs(t, err, ErrInvalidToken, "an expired service token must fail Parse")
}

func TestIssueServiceToken_RejectsBadParams(t *testing.T) {
	t.Parallel()
	s := testService("service-secret")

	cases := []struct {
		name    string
		subject string
		role    string
		ttl     time.Duration
	}{
		{"empty subject", "", "service", time.Hour},
		{"whitespace subject", "   ", "service", time.Hour},
		{"empty role", "service:risk", "", time.Hour},
		{"whitespace role", "service:risk", "\t", time.Hour},
		{"zero ttl", "service:risk", "service", 0},
		{"negative ttl", "service:risk", "service", -time.Minute},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := s.IssueServiceToken(tc.subject, tc.role, tc.ttl)
			assert.ErrorIs(t, err, ErrInvalidServiceToken)
		})
	}
}

func TestIssueServiceToken_DisabledService(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{}) // no secret

	_, err := s.IssueServiceToken("service:risk", "service", time.Hour)
	assert.ErrorIs(t, err, ErrTokensDisabled)
}
