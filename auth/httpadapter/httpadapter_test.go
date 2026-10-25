package httpadapter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ruoxizhnya/tellurion-kit/auth"
)

func TestDecideStartup(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		secret        string
		allowInsecure bool
		bindHost      string
		exposure      string
		wantSecret    bool
		wantInsecure  bool
		wantRefuse    bool
	}{
		{"secret enables auth", "s3cret", false, "0.0.0.0", "", true, false, false},
		{"secret wins over the insecure flag", "s3cret", true, "0.0.0.0", "", true, false, false},
		{"no secret and no opt-in refuses", "", false, "127.0.0.1", "", false, false, true},
		{"whitespace secret is not a secret", "   ", true, "0.0.0.0", "", false, false, true},
		{"no secret + loopback bind is open", "", true, "127.0.0.1", "", false, true, false},
		{"no secret + localhost is open", "", true, "localhost", "", false, true, false},
		{"no secret + non-loopback refuses", "", true, "0.0.0.0", "", false, false, true},
		{"container exemption is open", "", true, "0.0.0.0", InsecureExposureLoopbackPublished, false, true, false},
		{"misspelled exemption refuses", "", true, "0.0.0.0", "loopback", false, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := DecideStartup(tc.secret, tc.allowInsecure, tc.bindHost, tc.exposure)
			assert.Equal(t, tc.wantRefuse, got.Refuse)
			assert.Equal(t, tc.wantInsecure, got.Insecure)
			assert.Equal(t, tc.wantSecret, len(got.Secret) > 0)
			if got.Refuse {
				assert.NotEmpty(t, got.Reason, "a refusal must carry an actionable reason")
			}
		})
	}
}

func TestIsLoopbackHost(t *testing.T) {
	t.Parallel()
	cases := map[string]bool{
		"127.0.0.1": true,
		"::1":       true,
		"[::1]":     true,
		"localhost": true,
		"":          false,
		"0.0.0.0":   false,
		"10.0.0.1":  false,
		"garbage":   false,
	}
	for host, want := range cases {
		assert.Equal(t, want, IsLoopbackHost(host), host)
	}
}

func testPolicy() *Policy {
	return MustNewPolicy(map[string]Rule{
		"POST /api/auth/login":     {Public: true},
		"POST /api/widgets":        {Roles: []string{"trader", "admin"}},
		"DELETE /api/widgets/{id}": {Roles: []string{"trader", "admin"}},
		"POST /api/query":          {AnyAuthenticated: true},
	})
}

// newMux wires Guard onto every route, exactly as a service would
// (mux.Handle(pattern, Guard(...)(handler))). net/http sets r.Pattern from the
// registered pattern before the wrapped handler runs, which is what RouteKey
// reads.
func newMux(svc *auth.TokenService, p *Policy) *http.ServeMux {
	mux := http.NewServeMux()
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if claims, found := ClaimsFromContext(r.Context()); found {
			w.Header().Set("X-Role", claims.Role)
		}
		w.WriteHeader(http.StatusOK)
	})
	for _, pattern := range []string{
		"POST /api/auth/login",
		"POST /api/widgets",
		"DELETE /api/widgets/{id}",
		"POST /api/query",
		"GET /api/widgets",
		"POST /api/unclassified",
	} {
		mux.Handle(pattern, Guard(svc, p)(ok))
	}
	return mux
}

func newService(t *testing.T) *auth.TokenService {
	t.Helper()
	svc := auth.NewTokenService(auth.Config{JWTSecret: []byte("httpadapter-test")})
	require.True(t, svc.Enabled())
	return svc
}

func accessToken(t *testing.T, svc *auth.TokenService, role string) string {
	t.Helper()
	tok, _, err := svc.IssueTokens(1, role, role)
	require.NoError(t, err)
	return tok
}

func do(mux http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func TestGuard(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	mux := newMux(svc, testPolicy())
	trader := accessToken(t, svc, "trader")
	admin := accessToken(t, svc, "admin")
	viewer := accessToken(t, svc, "viewer")
	_, refresh, err := svc.IssueTokens(7, "refresher", "admin")
	require.NoError(t, err)

	cases := []struct {
		name   string
		method string
		path   string
		token  string
		want   int
	}{
		{"no token on a roles route", http.MethodPost, "/api/widgets", "", http.StatusUnauthorized},
		{"malformed token", http.MethodPost, "/api/widgets", "not-a-jwt", http.StatusUnauthorized},
		{"refresh token cannot be used as access", http.MethodPost, "/api/widgets", refresh, http.StatusUnauthorized},
		{"role too low", http.MethodPost, "/api/widgets", viewer, http.StatusForbidden},
		{"listed role passes", http.MethodPost, "/api/widgets", trader, http.StatusOK},
		{"second listed role passes", http.MethodPost, "/api/widgets", admin, http.StatusOK},
		{"templated route matches by pattern", http.MethodDelete, "/api/widgets/42", trader, http.StatusOK},
		{"any-authenticated viewer passes", http.MethodPost, "/api/query", viewer, http.StatusOK},
		{"any-authenticated needs a token", http.MethodPost, "/api/query", "", http.StatusUnauthorized},
		{"public route needs no token", http.MethodPost, "/api/auth/login", "", http.StatusOK},
		{"public route ignores a bad token", http.MethodPost, "/api/auth/login", "not-a-jwt", http.StatusOK},
		{"undeclared mutating route fails closed for admin", http.MethodPost, "/api/unclassified", admin, http.StatusForbidden},
		{"undeclared read route still needs a token", http.MethodGet, "/api/widgets", "", http.StatusUnauthorized},
		{"undeclared read route with a token passes", http.MethodGet, "/api/widgets", viewer, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, do(mux, tc.method, tc.path, tc.token).Code)
		})
	}
}

func TestGuard_UnknownMutatingRouteNamesTheRoute(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	mux := newMux(svc, testPolicy())

	w := do(mux, http.MethodPost, "/api/unclassified", accessToken(t, svc, "admin"))
	require.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "no authorization policy")
	assert.Contains(t, w.Body.String(), "/api/unclassified")
}

func TestGuard_AuthDisabledIsOpenAccess(t *testing.T) {
	t.Parallel()
	disabled := auth.NewTokenService(auth.Config{})
	require.False(t, disabled.Enabled())

	mux := newMux(disabled, testPolicy())
	assert.Equal(t, http.StatusOK, do(mux, http.MethodPost, "/api/widgets", "").Code)
	assert.Equal(t, http.StatusOK, do(mux, http.MethodPost, "/api/unclassified", "").Code,
		"disabled auth must not enforce missing policy entries either")

	var nilService *auth.TokenService
	assert.Equal(t, http.StatusOK, do(newMux(nilService, testPolicy()), http.MethodPost, "/api/widgets", "").Code)
}

func TestGuard_AttachesClaims(t *testing.T) {
	t.Parallel()
	svc := newService(t)
	mux := newMux(svc, testPolicy())

	w := do(mux, http.MethodPost, "/api/widgets", accessToken(t, svc, "trader"))
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "trader", w.Header().Get("X-Role"))
}

func TestAuthenticate_Sentinels(t *testing.T) {
	t.Parallel()
	svc := newService(t)

	_, err := Authenticate(svc, httptest.NewRequest(http.MethodGet, "/", nil))
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnauthenticated))
	assert.Contains(t, err.Error(), "missing Authorization header")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	_, err = Authenticate(svc, req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnauthenticated))

	_, refresh, err := svc.IssueTokens(3, "refresher", "admin")
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+refresh)
	_, err = Authenticate(svc, req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnauthenticated))
	assert.Contains(t, err.Error(), "refresh token")

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken(t, svc, "trader"))
	claims, err := Authenticate(svc, req)
	require.NoError(t, err)
	assert.Equal(t, "trader", claims.Role)

	var nilService *auth.TokenService
	_, err = Authenticate(nilService, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.True(t, errors.Is(err, ErrUnauthenticated))
}

func TestAuthorize(t *testing.T) {
	t.Parallel()
	viewer := &auth.Claims{Role: "viewer"}
	trader := &auth.Claims{Role: "trader"}

	cases := []struct {
		name    string
		claims  *auth.Claims
		rule    Rule
		wantErr error
	}{
		{"public needs no claims", nil, Rule{Public: true}, nil},
		{"any-authenticated without claims", nil, Rule{AnyAuthenticated: true}, ErrUnauthenticated},
		{"any-authenticated with claims", viewer, Rule{AnyAuthenticated: true}, nil},
		{"roles without claims", nil, Rule{Roles: []string{"trader"}}, ErrUnauthenticated},
		{"role not listed", viewer, Rule{Roles: []string{"trader", "admin"}}, ErrForbidden},
		{"role listed", trader, Rule{Roles: []string{"trader", "admin"}}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Authorize(tc.claims, tc.rule)
			if tc.wantErr == nil {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.True(t, errors.Is(err, tc.wantErr))
		})
	}
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
			t.Parallel()
			_, err := NewPolicy(map[string]Rule{"POST /api/x": rule})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "POST /api/x")
		})
	}

	_, err := NewPolicy(map[string]Rule{"": {Public: true}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestMustNewPolicy_PanicsOnMalformed(t *testing.T) {
	t.Parallel()
	assert.Panics(t, func() {
		MustNewPolicy(map[string]Rule{"POST /api/x": {}})
	})
}

func TestPolicy_ReturnsCopies(t *testing.T) {
	t.Parallel()
	p := testPolicy()

	rule, ok := p.Rule("POST /api/widgets")
	require.True(t, ok)
	rule.Roles[0] = "viewer"
	again, _ := p.Rule("POST /api/widgets")
	assert.Equal(t, []string{"trader", "admin"}, again.Roles)

	entries := p.Entries()
	entries["POST /api/widgets"] = Rule{Public: true}
	still, _ := p.Rule("POST /api/widgets")
	assert.False(t, still.Public)

	_, ok = p.Rule("POST /api/nope")
	assert.False(t, ok)
}

func TestRouteKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/api/x", nil)
	assert.Empty(t, RouteKey(req), "no resolved pattern means no key (fail closed)")

	req.Pattern = "POST /api/x/{id}"
	assert.Equal(t, "POST /api/x/{id}", RouteKey(req))

	req = httptest.NewRequest(http.MethodGet, "/api/x", nil)
	req.Pattern = "/api/x"
	assert.Equal(t, "GET /api/x", RouteKey(req))
}

func TestIsMutating(t *testing.T) {
	t.Parallel()
	for _, m := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		assert.True(t, IsMutating(m), m)
	}
	for _, m := range []string{http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace} {
		assert.False(t, IsMutating(m), m)
	}
}
