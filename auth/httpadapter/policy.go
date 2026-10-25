package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ruoxizhnya/tellurion-kit/auth"
)

// Rule is the access requirement for one route key. Exactly one tier must be
// selected; NewPolicy rejects anything else.
type Rule struct {
	Public           bool
	AnyAuthenticated bool
	Roles            []string
}

// Policy is an immutable route-key → Rule table. A route key is
// "METHOD /path/with/{param}" — the template the route was registered with,
// not the concrete URL.
type Policy struct {
	rules map[string]Rule
}

// NewPolicy validates and freezes the rule table. A malformed rule (two tiers
// at once, no tier at all, or an empty key) is a programming error and is
// reported here rather than silently resolving to allow or deny at request
// time.
func NewPolicy(rules map[string]Rule) (*Policy, error) {
	out := make(map[string]Rule, len(rules))
	for key, r := range rules {
		if strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("httpadapter: policy key must not be empty")
		}
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
			return nil, fmt.Errorf("httpadapter: rule %q must select exactly one tier (public / any-authenticated / roles), got %d", key, tiers)
		}
		cp := Rule{Public: r.Public, AnyAuthenticated: r.AnyAuthenticated}
		cp.Roles = append([]string(nil), r.Roles...)
		out[key] = cp
	}
	return &Policy{rules: out}, nil
}

// MustNewPolicy is NewPolicy for rule tables built as package-level values: it
// panics on a malformed rule so the mistake surfaces at init, not under load.
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
// HEAD/OPTIONS/TRACE are excluded for the same reason as in authz: they do
// not change state. The mutating surface is the four verbs below.
func IsMutating(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// RouteKey returns the policy key for r: the matched ServeMux pattern,
// prefixed with the method when the pattern does not already carry one.
// Returns "" when net/http did not resolve a pattern.
func RouteKey(r *http.Request) string {
	p := r.Pattern
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, r.Method+" ") {
		return p
	}
	return r.Method + " " + p
}

// Guard returns net/http middleware enforcing authentication and p.
//
// Semantics:
//   - auth disabled (svc == nil or !Enabled) → no-op, preserving the
//     development / CI open-access mode;
//   - a declared Public route → pass without a token;
//   - a mutating route with no policy entry → 403 naming the route
//     (fail-closed, even for an administrator);
//   - everything else needs a valid access token (401 otherwise), and a
//     declared Roles rule additionally needs one of the listed roles (403).
//
// On success the claims are attached to the request context; handlers read
// them with ClaimsFromContext(r.Context()).
func Guard(svc *auth.TokenService, p *Policy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if svc == nil || !svc.Enabled() {
				next.ServeHTTP(w, r)
				return
			}

			key := RouteKey(r)
			rule, declared := p.Rule(key)
			if declared && rule.Public {
				next.ServeHTTP(w, r)
				return
			}
			if !declared && IsMutating(r.Method) {
				writeJSON(w, http.StatusForbidden, map[string]any{
					"error": "no authorization policy for this endpoint",
					"route": key,
					"hint":  "declare this route in the service policy table (public / any-authenticated / roles)",
				})
				return
			}

			claims, reason, err := authenticate(svc, r)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]any{"error": reason})
				return
			}
			if declared {
				if err := Authorize(claims, rule); err != nil {
					if errors.Is(err, ErrForbidden) {
						writeJSON(w, http.StatusForbidden, map[string]any{
							"error": "insufficient role",
							"have":  claims.Role,
							"need":  rule.Roles,
						})
						return
					}
					writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "no authenticated user"})
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey{}, claims)))
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
