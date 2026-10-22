package authz

// Unit tests for the enforcement mechanism: tier semantics, fail-closed on an
// undeclared route, and the validation that rejects ambiguous rule entries.
// Both consuming services rely on these semantics, so they are pinned here
// once instead of being re-derived by each service's route-table guard test.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ruoxizhnya/tellurion-kit/auth"
)

func testPolicy() *Policy {
	return MustNewPolicy(map[string]Rule{
		"POST /api/auth/login":    {Public: true},
		"POST /api/admin/thing":   {Roles: []string{"admin"}},
		"POST /api/widgets":       {Roles: []string{"trader", "admin"}},
		"DELETE /api/widgets/:id": {Roles: []string{"trader", "admin"}},
		"POST /api/query":         {AnyAuthenticated: true},
	})
}

// testRouter builds a minimal auth + enforcement router. The second return
// value is unused by most tests but keeps registration and policy in one
// place.
func testRouter(svc *auth.TokenService, p *Policy) *gin.Engine {
	r := gin.New()
	if svc != nil {
		r.Use(svc.Middleware())
	}
	r.Use(Middleware(svc, p))
	register := func(method, path string) {
		r.Handle(method, path, func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	register(http.MethodPost, "/api/auth/login")
	register(http.MethodPost, "/api/admin/thing")
	register(http.MethodPost, "/api/widgets")
	register(http.MethodDelete, "/api/widgets/:id")
	register(http.MethodPost, "/api/query")
	register(http.MethodGet, "/api/widgets")
	// Deliberately registered but absent from the policy — validates
	// fail-closed enforcement.
	register(http.MethodPost, "/api/unclassified")
	return r
}

func enabledService(t *testing.T) *auth.TokenService {
	t.Helper()
	svc := auth.NewTokenService(auth.Config{
		JWTSecret:   []byte("authz-test"),
		PublicPaths: []string{"/api/auth/login"},
	})
	require.True(t, svc.Enabled())
	return svc
}

func tokenForRole(t *testing.T, svc *auth.TokenService, role string) string {
	t.Helper()
	tok, _, err := svc.IssueTokens(1, role, role)
	require.NoError(t, err)
	return tok
}

func do(r *gin.Engine, method, path, token string) int {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

// TestMiddleware_AuthDisabledIsNoOp: the development / CI open-access mode
// must be preserved, including for routes missing from the policy.
func TestMiddleware_AuthDisabledIsNoOp(t *testing.T) {
	t.Parallel()

	disabled := auth.NewTokenService(auth.Config{})
	require.False(t, disabled.Enabled())
	r := testRouter(disabled, testPolicy())

	assert.Equal(t, http.StatusOK, do(r, http.MethodPost, "/api/widgets", ""))
	assert.Equal(t, http.StatusOK, do(r, http.MethodPost, "/api/unclassified", ""),
		"disabled auth must not enforce missing policy entries either")

	// A nil service is the same posture (an unconfigured deployment).
	r2 := gin.New()
	r2.Use(Middleware(nil, testPolicy()))
	r2.POST("/api/widgets", func(c *gin.Context) { c.Status(http.StatusOK) })
	assert.Equal(t, http.StatusOK, do(r2, http.MethodPost, "/api/widgets", ""))
}

// TestMiddleware_RoleTiers pins the four tiers one by one.
func TestMiddleware_RoleTiers(t *testing.T) {
	t.Parallel()

	svc := enabledService(t)
	r := testRouter(svc, testPolicy())
	viewer := tokenForRole(t, svc, "viewer")
	trader := tokenForRole(t, svc, "trader")
	admin := tokenForRole(t, svc, "admin")

	cases := []struct {
		name   string
		method string
		path   string
		token  string
		want   int
	}{
		{"read endpoint needs no role (auth middleware already checked the token)", http.MethodGet, "/api/widgets", viewer, http.StatusOK},
		{"public mutating route passes without a token", http.MethodPost, "/api/auth/login", "", http.StatusOK},
		{"roles tier: viewer rejected", http.MethodPost, "/api/widgets", viewer, http.StatusForbidden},
		{"roles tier: trader passes", http.MethodPost, "/api/widgets", trader, http.StatusOK},
		{"roles tier: admin passes because the caller listed it", http.MethodPost, "/api/widgets", admin, http.StatusOK},
		{"path parameters match by template", http.MethodDelete, "/api/widgets/42", trader, http.StatusOK},
		{"admin tier: trader rejected", http.MethodPost, "/api/admin/thing", trader, http.StatusForbidden},
		{"admin tier: admin passes", http.MethodPost, "/api/admin/thing", admin, http.StatusOK},
		{"any-authenticated: viewer passes", http.MethodPost, "/api/query", viewer, http.StatusOK},
		{"any-authenticated: no token 401", http.MethodPost, "/api/query", "", http.StatusUnauthorized},
		{"roles tier: no token 401", http.MethodPost, "/api/widgets", "", http.StatusUnauthorized},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, do(r, tc.method, tc.path, tc.token))
		})
	}
}

// TestMiddleware_UnknownMutatingRouteFailsClosed: a missing policy entry is a
// hard 403 naming the route — never a silent pass, even for an admin.
func TestMiddleware_UnknownMutatingRouteFailsClosed(t *testing.T) {
	t.Parallel()

	svc := enabledService(t)
	r := testRouter(svc, testPolicy())
	admin := tokenForRole(t, svc, "admin")

	req := httptest.NewRequest(http.MethodPost, "/api/unclassified", nil)
	req.Header.Set("Authorization", "Bearer "+admin)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code, "even an admin must not pass an undeclared endpoint")
	assert.Contains(t, w.Body.String(), "no authorization policy")
	assert.Contains(t, w.Body.String(), "/api/unclassified")
}

// TestMiddleware_NilPolicyFailsClosed: no table at all must not mean "allow".
func TestMiddleware_NilPolicyFailsClosed(t *testing.T) {
	t.Parallel()
	svc := enabledService(t)
	r := testRouter(svc, nil)
	trader := tokenForRole(t, svc, "trader")
	assert.Equal(t, http.StatusForbidden, do(r, http.MethodPost, "/api/widgets", trader))
}

// TestMiddleware_ReadMethodsBypassPolicy proves this layer governs only the
// mutating surface: with enforcement on and no policy at all, reads still
// pass and writes still fail closed. Uses a stub Authenticator so the result
// cannot come from auth middleware.
func TestMiddleware_ReadMethodsBypassPolicy(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.Use(Middleware(stubAuth{enabled: true}, nil))
	r.GET("/api/anything", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/api/anything", func(c *gin.Context) { c.Status(http.StatusOK) })

	assert.Equal(t, http.StatusOK, do(r, http.MethodGet, "/api/anything", ""))
	assert.Equal(t, http.StatusForbidden, do(r, http.MethodPost, "/api/anything", ""))
}

type stubAuth struct{ enabled bool }

func (s stubAuth) Enabled() bool { return s.enabled }

func TestPolicy_Rule(t *testing.T) {
	t.Parallel()

	p := testPolicy()

	rule, ok := p.Rule("POST /api/auth/login")
	require.True(t, ok)
	assert.True(t, rule.Public)

	rule, ok = p.Rule("POST /api/admin/thing")
	require.True(t, ok)
	assert.Equal(t, []string{"admin"}, rule.Roles)

	rule, ok = p.Rule("POST /api/query")
	require.True(t, ok)
	assert.True(t, rule.AnyAuthenticated)

	_, ok = p.Rule("POST /api/nope")
	assert.False(t, ok)

	assert.Len(t, p.Entries(), 5)
}

// TestPolicy_ReturnsCopies: the table is immutable once built; mutating a
// returned Rule or Entries map must not affect enforcement.
func TestPolicy_ReturnsCopies(t *testing.T) {
	t.Parallel()
	p := testPolicy()

	rule, _ := p.Rule("POST /api/widgets")
	rule.Roles[0] = "viewer"
	again, _ := p.Rule("POST /api/widgets")
	assert.Equal(t, []string{"trader", "admin"}, again.Roles)

	entries := p.Entries()
	entries["POST /api/widgets"] = Rule{Public: true}
	still, _ := p.Rule("POST /api/widgets")
	assert.False(t, still.Public)
}

func TestNewPolicy_RejectsMalformedRules(t *testing.T) {
	t.Parallel()
	cases := map[string]Rule{
		"two tiers at once": {Public: true, Roles: []string{"admin"}},
		"public and any":    {Public: true, AnyAuthenticated: true},
		"any and roles":     {AnyAuthenticated: true, Roles: []string{"admin"}},
		"no tier selected":  {},
	}
	for name, rule := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := NewPolicy(map[string]Rule{"POST /api/x": rule})
			require.Error(t, err, "a rule must select exactly one tier")
			assert.Contains(t, err.Error(), "POST /api/x")
		})
	}
}

func TestMustNewPolicy_PanicsOnMalformed(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		MustNewPolicy(map[string]Rule{"POST /api/x": {}})
	})
}

func TestIsMutating(t *testing.T) {
	t.Parallel()
	for _, m := range []string{"POST", "PUT", "PATCH", "DELETE"} {
		assert.True(t, IsMutating(m), m)
	}
	for _, m := range []string{"GET", "HEAD", "OPTIONS", "TRACE"} {
		assert.False(t, IsMutating(m), m)
	}
}
