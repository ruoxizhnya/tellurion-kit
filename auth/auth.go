// Package auth provides the mechanism half of token-based authentication:
// HS256 JWT issuance / parsing / refresh, bcrypt password hashing, and gin
// middleware that puts the caller identity on the request context.
//
// It deliberately holds no storage. User rows, password verification against
// a database, audit persistence and the role vocabulary live in the service
// that owns them (SCOPE.md §5 deviation #3): this package never opens a
// connection, never names a table, and treats roles as opaque strings.
//
// Disabled mode: when Config.JWTSecret is empty the service is inert —
// Middleware and RequireRole become no-ops and token operations fail. This is
// the documented "open access" posture for development and CI; it is a single
// switch so a deployment cannot end up half-protected.
package auth

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// BcryptCost is the work factor used for new password hashes. 12 is the
// 2026 industry default — slow enough to discourage offline cracking on
// consumer GPUs, fast enough to not tank login latency.
const BcryptCost = 12

// Sentinel errors. Callers compare with errors.Is.
var (
	ErrInvalidToken        = errors.New("auth: invalid or expired token")
	ErrUserDisabled        = errors.New("auth: user account is disabled")
	ErrTokensDisabled      = errors.New("auth: cannot issue tokens: JWT secret not configured")
	ErrUnauthenticated     = errors.New("auth: request has no authenticated user")
	ErrMissingIdentity     = errors.New("auth: refresh requires an identity loader")
	ErrIdentityLoadFailed  = errors.New("auth: identity lookup failed")
	ErrInvalidServiceToken = errors.New("auth: invalid service token parameters")
)

// Config is the static configuration for token issuance.
type Config struct {
	JWTSecret       []byte
	AccessTokenTTL  time.Duration // default 15m
	RefreshTokenTTL time.Duration // default 7d
	Issuer          string        // default "tellurion"
	// PublicPaths is the caller-supplied whitelist of request paths that skip
	// the JWT middleware entirely (login, refresh, health, metrics, …).
	// The route table itself belongs to the application; this package only
	// enforces whatever list it is handed.
	PublicPaths []string
}

// WithDefaults returns c with zero-value fields replaced by sensible defaults.
func (c Config) WithDefaults() Config {
	if c.AccessTokenTTL == 0 {
		c.AccessTokenTTL = 15 * time.Minute
	}
	if c.RefreshTokenTTL == 0 {
		c.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if c.Issuer == "" {
		c.Issuer = "tellurion"
	}
	return c
}

// TokenService mints and verifies tokens. It is safe for concurrent use.
// Build one per process; it holds no per-request state.
type TokenService struct {
	cfg         Config
	secret      []byte
	publicPaths map[string]struct{}
}

// NewTokenService creates a service. An empty Config.JWTSecret yields a
// disabled service (see the package comment).
func NewTokenService(cfg Config) *TokenService {
	cfg = cfg.WithDefaults()
	paths := make(map[string]struct{}, len(cfg.PublicPaths))
	for _, p := range cfg.PublicPaths {
		paths[p] = struct{}{}
	}
	return &TokenService{
		cfg:         cfg,
		secret:      cfg.JWTSecret,
		publicPaths: paths,
	}
}

// Enabled reports whether JWT issuance / verification is on.
func (s *TokenService) Enabled() bool {
	return s != nil && len(s.secret) > 0
}

// AccessTTL returns the configured access token TTL. Handlers use it to
// populate the expires_in field of a token response.
func (s *TokenService) AccessTTL() time.Duration {
	return s.cfg.AccessTokenTTL
}

// PublicPaths returns a copy of the unauthenticated path whitelist. Exported
// so a route table can be machine-checked against it.
func (s *TokenService) PublicPaths() []string {
	out := make([]string, 0, len(s.publicPaths))
	for p := range s.publicPaths {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// IsPublicPath reports whether path skips authentication.
func (s *TokenService) IsPublicPath(path string) bool {
	_, ok := s.publicPaths[path]
	return ok
}

// Identity is the token-relevant view of the authenticated principal. It is
// an input shape, not a row: the caller maps its own user record onto it.
type Identity struct {
	ID       int64
	Username string
	Role     string
	Disabled bool
}

// IdentityLoader looks up the current state of the principal identified by
// the refresh token's subject. Refresh needs it so a token cannot outlive a
// disabled account.
type IdentityLoader func(ctx context.Context, id int64) (Identity, error)

// IssueTokens mints a fresh access + refresh token pair.
func (s *TokenService) IssueTokens(id int64, username, role string) (access, refresh string, err error) {
	if !s.Enabled() {
		return "", "", ErrTokensDisabled
	}
	now := time.Now().UTC()

	accessClaims := Claims{
		Username: username,
		Role:     role,
		Kind:     KindAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   formatUserID(id),
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
		},
	}
	refreshClaims := Claims{
		Username: username,
		Role:     role,
		Kind:     KindRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   formatUserID(id),
			Issuer:    s.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.RefreshTokenTTL)),
		},
	}

	access, err = s.mint(accessClaims)
	if err != nil {
		return "", "", err
	}
	refresh, err = s.mint(refreshClaims)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ParseRefresh validates a refresh token and returns its claims. Any failure
// (bad signature, expiry, wrong kind) yields ErrInvalidToken.
func (s *TokenService) ParseRefresh(token string) (*Claims, error) {
	claims, err := s.Parse(token)
	if err != nil {
		return nil, err
	}
	if claims.Kind != KindRefresh {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// Refresh validates a refresh token, re-reads the principal through load (so
// a disabled account cannot keep rotating tokens), and mints a fresh pair.
// Any failure that is not ErrUserDisabled is reported as ErrInvalidToken to
// keep token-state details out of the response.
func (s *TokenService) Refresh(ctx context.Context, refreshToken string, load IdentityLoader) (access, refresh string, err error) {
	if load == nil {
		return "", "", ErrMissingIdentity
	}
	claims, err := s.ParseRefresh(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}
	uid, err := claims.UserIDInt64()
	if err != nil {
		return "", "", ErrInvalidToken
	}
	id, err := load(ctx, uid)
	if err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrIdentityLoadFailed, err)
	}
	if id.Disabled {
		return "", "", ErrUserDisabled
	}
	return s.IssueTokens(id.ID, id.Username, id.Role)
}

func (s *TokenService) mint(c Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(s.secret)
}

// Parse validates the token's signature and expiry and returns the claims.
func (s *TokenService) Parse(token string) (*Claims, error) {
	if !s.Enabled() {
		return nil, ErrInvalidToken
	}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}))
	var c Claims
	t, err := parser.ParseWithClaims(token, &c, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil || !t.Valid {
		return nil, ErrInvalidToken
	}
	return &c, nil
}

// HashPassword returns a bcrypt hash at BcryptCost. The plaintext is never
// retained or logged.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword reports whether password matches hash. It returns nil on a
// match and a non-nil error otherwise; callers are expected to collapse the
// error into their own "invalid credentials" response so account-existence
// is not leaked.
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
