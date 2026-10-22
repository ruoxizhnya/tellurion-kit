package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func testService(secret string) *TokenService {
	return NewTokenService(Config{
		JWTSecret:   []byte(secret),
		PublicPaths: []string{"/api/auth/login", "/health"},
	})
}

func TestConfig_WithDefaults(t *testing.T) {
	t.Parallel()
	c := Config{}.WithDefaults()
	assert.Equal(t, 15*time.Minute, c.AccessTokenTTL)
	assert.Equal(t, 7*24*time.Hour, c.RefreshTokenTTL)
	assert.Equal(t, "tellurion", c.Issuer)

	c2 := Config{AccessTokenTTL: time.Minute, Issuer: "foo"}.WithDefaults()
	assert.Equal(t, time.Minute, c2.AccessTokenTTL)
	assert.Equal(t, 7*24*time.Hour, c2.RefreshTokenTTL)
	assert.Equal(t, "foo", c2.Issuer)
}

func TestTokenService_Enabled(t *testing.T) {
	t.Parallel()
	assert.True(t, testService("hello").Enabled())
	assert.False(t, NewTokenService(Config{}).Enabled())

	var nilSvc *TokenService
	assert.False(t, nilSvc.Enabled(), "nil receiver must read as disabled, not panic")
}

func TestTokenService_Disabled_IssuesAndParsesFail(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{}) // no secret
	_, _, err := s.IssueTokens(1, "alice", "admin")
	assert.ErrorIs(t, err, ErrTokensDisabled)

	_, err = s.Parse("anything")
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenService_IssueAndParse_RoundTrip(t *testing.T) {
	t.Parallel()
	s := testService("test-secret")
	access, refresh, err := s.IssueTokens(42, "alice", "trader")
	require.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	assert.NotEqual(t, access, refresh)

	claims, err := s.Parse(access)
	require.NoError(t, err)
	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "trader", claims.Role)
	assert.Equal(t, KindAccess, claims.Kind)
	assert.Equal(t, "42", claims.Subject)
	assert.Equal(t, "tellurion", claims.Issuer)
	uid, err := claims.UserIDInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(42), uid)

	claims2, err := s.ParseRefresh(refresh)
	require.NoError(t, err)
	assert.Equal(t, KindRefresh, claims2.Kind)
}

func TestTokenService_Parse_RejectsBadSignature(t *testing.T) {
	t.Parallel()
	s1 := testService("secret-a")
	s2 := testService("secret-b")
	access, _, err := s1.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	_, err = s2.Parse(access)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

// TestTokenService_Parse_RejectsNonHMAC pins the algorithm allow-list: a
// token asking to be verified with alg=none (or any non-HS256 method) must
// not be accepted, regardless of what the keyfunc would return.
func TestTokenService_Parse_RejectsNonHMAC(t *testing.T) {
	t.Parallel()
	s := testService("test-secret")
	none := jwt.NewWithClaims(jwt.SigningMethodNone, Claims{
		Username: "mallory",
		Role:     "admin",
		Kind:     KindAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	signed, err := none.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = s.Parse(signed)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenService_Refresh_RotatesPair(t *testing.T) {
	t.Parallel()
	s := testService("test")
	_, refresh, err := s.IssueTokens(7, "alice", "trader")
	require.NoError(t, err)

	var askedFor int64
	load := func(_ context.Context, id int64) (Identity, error) {
		askedFor = id
		return Identity{ID: id, Username: "alice", Role: "trader"}, nil
	}

	access2, refresh2, err := s.Refresh(context.Background(), refresh, load)
	require.NoError(t, err)
	assert.Equal(t, int64(7), askedFor, "the loader must be asked for the subject in the refresh token")
	assert.NotEmpty(t, access2)
	assert.NotEmpty(t, refresh2)

	claims, err := s.Parse(access2)
	require.NoError(t, err)
	assert.Equal(t, KindAccess, claims.Kind)
	assert.Equal(t, "alice", claims.Username)

	claims, err = s.ParseRefresh(refresh2)
	require.NoError(t, err)
	assert.Equal(t, KindRefresh, claims.Kind)
}

func TestTokenService_Refresh_RejectsAccessToken(t *testing.T) {
	t.Parallel()
	s := testService("test")
	access, _, err := s.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	_, _, err = s.Refresh(context.Background(), access, func(_ context.Context, id int64) (Identity, error) {
		t.Fatal("loader must not be called for an access token")
		return Identity{}, nil
	})
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenService_Refresh_DisabledIdentityFails(t *testing.T) {
	t.Parallel()
	s := testService("test")
	_, refresh, err := s.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	_, _, err = s.Refresh(context.Background(), refresh, func(_ context.Context, id int64) (Identity, error) {
		return Identity{ID: id, Disabled: true}, nil
	})
	assert.ErrorIs(t, err, ErrUserDisabled)
}

func TestTokenService_Refresh_LoaderErrorIsWrapped(t *testing.T) {
	t.Parallel()
	s := testService("test")
	_, refresh, err := s.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	dbDown := errors.New("connection refused")
	_, _, err = s.Refresh(context.Background(), refresh, func(_ context.Context, _ int64) (Identity, error) {
		return Identity{}, dbDown
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrIdentityLoadFailed)
	assert.ErrorIs(t, err, dbDown, "the underlying cause must stay inspectable")
}

func TestTokenService_Refresh_NilLoaderFailsLoud(t *testing.T) {
	t.Parallel()
	s := testService("test")
	_, refresh, err := s.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	_, _, err = s.Refresh(context.Background(), refresh, nil)
	assert.ErrorIs(t, err, ErrMissingIdentity)
}

func TestTokenService_Refresh_GarbageFails(t *testing.T) {
	t.Parallel()
	s := testService("test")
	_, _, err := s.Refresh(context.Background(), "not-a-jwt", func(_ context.Context, id int64) (Identity, error) {
		return Identity{ID: id}, nil
	})
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestClaims_UserIDInt64_Errors(t *testing.T) {
	t.Parallel()
	c := &Claims{}
	_, err := c.UserIDInt64()
	assert.Error(t, err)

	c.Subject = "not-a-number"
	_, err = c.UserIDInt64()
	assert.Error(t, err)
}

func TestHashPassword_RoundTrip(t *testing.T) {
	t.Parallel()
	hash, err := HashPassword("s3cret-password")
	require.NoError(t, err)
	assert.NotEqual(t, "s3cret-password", hash)
	assert.False(t, strings.Contains(hash, "s3cret"), "hash must not embed the plaintext")

	cost, err := bcrypt.Cost([]byte(hash))
	require.NoError(t, err)
	assert.Equal(t, BcryptCost, cost, "new hashes must use the package work factor")

	assert.NoError(t, VerifyPassword(hash, "s3cret-password"))
	assert.Error(t, VerifyPassword(hash, "wrong"))
}

func TestPublicPaths_CopyAndLookup(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{
		JWTSecret:   []byte("x"),
		PublicPaths: []string{"/health", "/api/auth/login"},
	})
	assert.True(t, s.IsPublicPath("/health"))
	assert.True(t, s.IsPublicPath("/api/auth/login"))
	assert.False(t, s.IsPublicPath("/api/private"))

	got := s.PublicPaths()
	assert.Equal(t, []string{"/api/auth/login", "/health"}, got)
	got[0] = "/mutated"
	assert.True(t, s.IsPublicPath("/api/auth/login"), "PublicPaths must return a copy")
}

func TestAccessTTL(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{JWTSecret: []byte("x"), AccessTokenTTL: 3 * time.Minute})
	assert.Equal(t, 3*time.Minute, s.AccessTTL())
}
