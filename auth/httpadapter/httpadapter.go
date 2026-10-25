// Package httpadapter applies tellurion-kit/auth to net/http services.
//
// It is the framework-independent counterpart of the gin middleware in
// auth/authz: the same Bearer-token parsing, the same "refresh token cannot
// be used as an access token" rule, and the same any-of role check. It
// deliberately imports no gin, so net/http services can use it without
// pulling the gin stack in.
//
// Enforcement shape: one Policy per service maps route keys
// ("METHOD /path/with/{param}") to Public / AnyAuthenticated / Roles. A
// mutating request whose route is missing from the policy is rejected with
// 403 — never silently allowed. Route keys come from http.Request.Pattern,
// which net/http fills in from the registered ServeMux pattern before calling
// the handler, so install Guard on each route:
//
//	mux.Handle("POST /api/widgets", httpadapter.Guard(svc, policy)(h))
//
// A request that reaches Guard without a resolved pattern (for example when
// Guard wraps the whole mux instead of a route) cannot be looked up: mutating
// requests fail closed with 403, reads with 401.
package httpadapter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ruoxizhnya/tellurion-kit/auth"
)

// Sentinel errors. Callers compare with errors.Is; Guard maps them to 401
// and 403 respectively.
var (
	ErrUnauthenticated = errors.New("httpadapter: unauthenticated")
	ErrForbidden       = errors.New("httpadapter: forbidden")
)

// InsecureExposureLoopbackPublished is the only accepted exposure
// declaration that permits open access on a non-loopback bind: the process
// binds 0.0.0.0 (as containers must) while the publishing layer maps the port
// to 127.0.0.1 only. Anything else — misspelled, empty, partial — counts as
// "not declared" and the service refuses to start.
const InsecureExposureLoopbackPublished = "loopback-published"

// Startup is the outcome of the startup gate. It is a pure value so the
// decision is testable; the caller performs the actual exit/fatal.
type Startup struct {
	Secret   []byte // non-empty = authentication enabled
	Insecure bool   // explicitly running without authentication
	Refuse   bool   // refuse to start
	Reason   string // actionable explanation when Refuse is true
}

// DecideStartup decides whether a service may start with the given auth
// configuration. Semantics match core's bootstrap.DecideAuthStartup:
//
//   - a non-empty secret turns authentication on;
//   - no secret + allowInsecure=false → refuse (never start wide open);
//   - no secret + allowInsecure=true + loopback bind → open access;
//   - no secret + allowInsecure=true + exposure=loopback-published → open
//     access (container form: bind 0.0.0.0, publish on loopback);
//   - anything else → refuse.
func DecideStartup(secret string, allowInsecure bool, bindHost, exposure string) Startup {
	if s := strings.TrimSpace(secret); s != "" {
		return Startup{Secret: []byte(s)}
	}
	if !allowInsecure {
		return Startup{
			Refuse: true,
			Reason: "auth: JWT secret missing — refusing to start in open-access mode. " +
				"Set JWT_SECRET (or the service's jwt_secret setting) to a strong random value. " +
				"Local dev only: allow insecure mode AND bind 127.0.0.1",
		}
	}
	if IsLoopbackHost(bindHost) {
		return Startup{Insecure: true}
	}
	if strings.TrimSpace(exposure) == InsecureExposureLoopbackPublished {
		return Startup{Insecure: true}
	}
	return Startup{
		Refuse: true,
		Reason: fmt.Sprintf("auth: insecure mode requested but bind host %q is not loopback — "+
			"open access would be reachable from other hosts. Bind 127.0.0.1, or set JWT_SECRET. "+
			"Containers that must bind 0.0.0.0 have to publish on loopback and declare exposure=%s",
			bindHost, InsecureExposureLoopbackPublished),
	}
}

// IsLoopbackHost reports whether host only listens on this machine. An empty
// host counts as non-loopback: net/http binding ":port" is 0.0.0.0.
func IsLoopbackHost(host string) bool {
	h := strings.TrimSpace(host)
	if h == "" {
		return false
	}
	if h == "localhost" {
		return true
	}
	ip := net.ParseIP(strings.Trim(h, "[]"))
	return ip != nil && ip.IsLoopback()
}

// Authenticate validates the Bearer access token on r and returns its claims.
// Every failure wraps ErrUnauthenticated, so callers can map it to 401:
// missing header, non-Bearer scheme, bad signature or expiry, a refresh token
// used as an access token, or a malformed subject.
func Authenticate(svc *auth.TokenService, r *http.Request) (*auth.Claims, error) {
	claims, reason, err := authenticate(svc, r)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthenticated, reason)
	}
	return claims, nil
}

// authenticate is Authenticate plus a client-safe reason for the 401 body.
func authenticate(svc *auth.TokenService, r *http.Request) (*auth.Claims, string, error) {
	if svc == nil || !svc.Enabled() {
		return nil, "authentication is not configured", ErrUnauthenticated
	}
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, "missing Authorization header", ErrUnauthenticated
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, "Authorization header must be 'Bearer <token>'", ErrUnauthenticated
	}
	claims, err := svc.Parse(parts[1])
	if err != nil {
		return nil, "invalid or expired token", ErrUnauthenticated
	}
	if claims.Kind != auth.KindAccess {
		return nil, "refresh token cannot be used as access token", ErrUnauthenticated
	}
	if _, err := claims.UserIDInt64(); err != nil {
		return nil, "malformed token subject", ErrUnauthenticated
	}
	return claims, "", nil
}

// Authorize applies rule to an already-authenticated principal. A Public rule
// needs no claims; AnyAuthenticated needs a principal; Roles needs one of the
// listed roles (any-of — callers express hierarchy by listing every
// acceptable role; this package does not rank roles).
func Authorize(claims *auth.Claims, rule Rule) error {
	if rule.Public {
		return nil
	}
	if claims == nil {
		return fmt.Errorf("%w: no authenticated user", ErrUnauthenticated)
	}
	if rule.AnyAuthenticated {
		return nil
	}
	for _, allowed := range rule.Roles {
		if allowed == claims.Role {
			return nil
		}
	}
	return fmt.Errorf("%w: role %q is not in %v", ErrForbidden, claims.Role, rule.Roles)
}

type claimsKey struct{}

// ClaimsFromContext returns the claims that Guard attached to the request.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*auth.Claims)
	return claims, ok
}
