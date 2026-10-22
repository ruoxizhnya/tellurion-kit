// Package authz turns "endpoint × required role" into an explicit table and
// enforces it. It supplies the mechanism only: the table itself — which route
// needs which role — is the application's asset and is passed in (SCOPE.md §3,
// deviation #4). Role names are opaque strings; this package knows nothing
// about any particular role vocabulary or hierarchy.
//
// Enforcement is table-driven rather than per-route middleware registration:
// middleware attached route by route duplicates the route table, and a missed
// attachment is invisible. Instead each service keeps one Policy covering all
// mutating endpoints, and a guard test can assert "route table ≡ policy".
//
// The failure mode is deliberately fail-closed: a mutating request whose route
// is missing from the policy is rejected with 403, even for an administrator.
// A guard test should have caught the omission in CI; this is the runtime
// backstop.
package authz

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ruoxizhnya/tellurion-kit/auth"
)

// Rule is the access requirement for one route key.
//
// Exactly one tier must be selected, enforced by NewPolicy:
//   - Public: no authentication at all (login endpoints);
//   - AnyAuthenticated: any valid token, regardless of role;
//   - Roles: any-of the listed roles (the caller expresses hierarchy by
//     listing every acceptable role — this package does not rank roles).
type Rule struct {
	Public           bool
	AnyAuthenticated bool
	Roles            []string
}

// Policy is a route-key → Rule table. A route key is the full route template
// in gin's form: "METHOD /path/with/:param".
type Policy struct {
	rules map[string]Rule
}

// NewPolicy validates and freezes the rule table. A malformed rule (two tiers
// at once, or no tier at all) is a programming error and is reported here
// rather than silently resolving to allow or deny at request time.
func NewPolicy(rules map[string]Rule) (*Policy, error) {
	out := make(map[string]Rule, len(rules))
	for key, r := range rules {
		tiers := 0
		if r.Public {
			tiers++
		}
		if r.AnyAuthenticated {
			tiers++
		}
		if len(r.Roles) > 0 {
			tiers++
		}
		if tiers != 1 {
			return nil, fmt.Errorf("authz: rule %q must select exactly one tier (public / any-authenticated / roles), got %d", key, tiers)
		}
		cp := Rule{Public: r.Public, AnyAuthenticated: r.AnyAuthenticated}
		cp.Roles = append([]string(nil), r.Roles...)
		out[key] = cp
	}
	return &Policy{rules: out}, nil
}

// MustNewPolicy is NewPolicy for rule tables built as package-level values:
// it panics on a malformed rule so the mistake surfaces at init, not under
// load. Keep it to startup wiring.
func MustNewPolicy(rules map[string]Rule) *Policy {
	p, err := NewPolicy(rules)
	if err != nil {
		panic(err)
	}
	return p
}

// Rule returns the requirement for a route key and whether it is registered.
func (p *Policy) Rule(key string) (Rule, bool) {
	if p == nil {
		return Rule{}, false
	}
	r, ok := p.rules[key]
	if !ok {
		return Rule{}, false
	}
	cp := r
	cp.Roles = append([]string(nil), r.Roles...)
	return cp, true
}

// Entries returns a copy of the whole table, keyed by route, for guard tests
// that reconcile it against the router's registered routes.
func (p *Policy) Entries() map[string]Rule {
	if p == nil {
		return nil
	}
	out := make(map[string]Rule, len(p.rules))
	for k, r := range p.rules {
		cp := r
		cp.Roles = append([]string(nil), r.Roles...)
		out[k] = cp
	}
	return out
}

// IsMutating reports whether method changes state.
//
// HEAD/OPTIONS are excluded (gin's Any registers them, but they do not change
// state); TRACE is excluded for the same reason. The mutating surface is the
// four verbs below.
func IsMutating(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// Authenticator is the minimal view of an auth service: whether enforcement
// is on. Both auth.TokenService and a service-local wrapper satisfy it.
type Authenticator interface {
	Enabled() bool
}

// Middleware enforces p. It MUST run after auth.Middleware() — it reads the
// role that middleware writes into the context.
//
// Semantics:
//   - auth disabled (svc == nil or !Enabled) → no-op, preserving the
//     development / CI open-access mode;
//   - read methods (GET/HEAD/OPTIONS…) → pass; this package only governs the
//     mutating surface;
//   - mutating methods → look up the route template: Public passes,
//     AnyAuthenticated needs a token, Roles needs one of the listed roles;
//     a missing entry is 403 with the route named (fail-closed).
func Middleware(svc Authenticator, p *Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc == nil || !svc.Enabled() {
			c.Next()
			return
		}
		if !IsMutating(c.Request.Method) {
			c.Next()
			return
		}

		// Use the route template, not the raw URL: the latter carries real
		// ids and cannot be looked up. An unmatched route has an empty
		// FullPath(); gin's 404 handles it, no need to intervene.
		fullPath := c.FullPath()
		if fullPath == "" {
			c.Next()
			return
		}
		key := c.Request.Method + " " + fullPath

		rule, ok := p.Rule(key)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no authorization policy for this endpoint",
				"route": key,
				"hint":  "declare this route in the service policy table (public / any-authenticated / roles)",
			})
			return
		}
		if rule.Public {
			c.Next()
			return
		}

		role, ok := auth.RoleFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authenticated user (authz.Middleware placed before auth.Middleware?)",
			})
			return
		}
		if rule.AnyAuthenticated {
			c.Next()
			return
		}
		for _, allowed := range rule.Roles {
			if allowed == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient role",
			"have":  role,
			"need":  rule.Roles,
		})
	}
}
