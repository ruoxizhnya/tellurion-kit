package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(s *TokenService) *gin.Engine {
	r := gin.New()
	r.Use(s.Middleware())
	r.GET("/api/secret", func(c *gin.Context) {
		uid, role, ok := UserFromContext(c)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"user_id": int64(0), "role": "anonymous"})
			return
		}
		name, _ := UsernameFromContext(c)
		_, hasClaims := ClaimsFromContext(c)
		c.JSON(http.StatusOK, gin.H{
			"user_id":    uid,
			"role":       role,
			"username":   name,
			"has_claims": hasClaims,
		})
	})
	r.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"login": "ok"})
	})
	return r
}

func TestMiddleware_Disabled_NoOp(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{}) // no secret
	r := setupRouter(s)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/secret", nil))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "anonymous")
}

func TestMiddleware_NilService_NoOp(t *testing.T) {
	t.Parallel()
	var s *TokenService
	r := setupRouter(s)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/secret", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_NoToken_401(t *testing.T) {
	t.Parallel()
	r := setupRouter(testService("test"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/api/secret", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing Authorization header")
}

func TestMiddleware_BadFormat_401(t *testing.T) {
	t.Parallel()
	r := setupRouter(testService("test"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/secret", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMiddleware_ValidToken_OK(t *testing.T) {
	t.Parallel()
	s := testService("test")
	r := setupRouter(s)
	access, _, err := s.IssueTokens(99, "alice", "trader")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/secret", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"user_id":99`)
	assert.Contains(t, w.Body.String(), `"role":"trader"`)
	assert.Contains(t, w.Body.String(), `"username":"alice"`)
	assert.Contains(t, w.Body.String(), `"has_claims":true`)
}

func TestMiddleware_RefreshTokenRejected(t *testing.T) {
	t.Parallel()
	s := testService("test")
	r := setupRouter(s)
	_, refresh, err := s.IssueTokens(1, "x", "admin")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/secret", nil)
	req.Header.Set("Authorization", "Bearer "+refresh)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "refresh token cannot be used as access token")
}

func TestMiddleware_PublicPath_NoAuth(t *testing.T) {
	t.Parallel()
	r := setupRouter(testService("test"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMiddleware_UnlistedPathStillGuarded(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{JWTSecret: []byte("test")}) // empty whitelist
	r := setupRouter(s)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/api/auth/login", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code,
		"paths not handed to Config.PublicPaths must not be treated as public")
}

func TestMiddleware_ExpiredToken_401(t *testing.T) {
	t.Parallel()
	s := testService("test")
	r := setupRouter(s)
	c := Claims{
		Username: "alice",
		Role:     "admin",
		Kind:     KindAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			Issuer:    "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte("test"))
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/secret", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRole_403(t *testing.T) {
	t.Parallel()
	s := testService("test")
	r := gin.New()
	r.Use(s.Middleware())
	r.GET("/admin", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	access, _, err := s.IssueTokens(1, "viewer", "viewer")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient role")
}

func TestRequireRole_OK(t *testing.T) {
	t.Parallel()
	s := testService("test")
	r := gin.New()
	r.Use(s.Middleware())
	r.GET("/admin", RequireRole("admin", "trader"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	access, _, err := s.IssueTokens(1, "trader", "trader")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_NoUser_401(t *testing.T) {
	t.Parallel()
	r := gin.New()
	r.GET("/admin", RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestServiceRequireRole_Disabled_AllowsEverything pins the short-circuit
// that makes TokenService.RequireRole usable behind TokenService.Middleware().
// If someone "simplifies" a call site back to the package-level RequireRole,
// the disabled-mode request below turns 401 and this test goes red.
func TestServiceRequireRole_Disabled_AllowsEverything(t *testing.T) {
	t.Parallel()
	s := NewTokenService(Config{})
	require.False(t, s.Enabled())

	r := gin.New()
	r.Use(s.Middleware())
	r.POST("/widgets", s.RequireRole("trader", "admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/widgets", nil))
	assert.Equal(t, http.StatusOK, w.Code,
		"auth disabled must stay open: RequireRole has to no-op when Enabled()==false")
}

// TestServiceRequireRole_Enabled_StillEnforces is the counter-test:
// the short-circuit must not turn into "always allow".
func TestServiceRequireRole_Enabled_StillEnforces(t *testing.T) {
	t.Parallel()
	s := testService("test")
	require.True(t, s.Enabled())

	r := gin.New()
	r.Use(s.Middleware())
	r.POST("/widgets", s.RequireRole("trader", "admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	viewer, _, err := s.IssueTokens(1, "v", "viewer")
	require.NoError(t, err)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/widgets", nil)
	req.Header.Set("Authorization", "Bearer "+viewer)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	trader, _, err := s.IssueTokens(2, "t", "trader")
	require.NoError(t, err)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/widgets", nil)
	req.Header.Set("Authorization", "Bearer "+trader)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("POST", "/widgets", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
