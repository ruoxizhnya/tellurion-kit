package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context keys for storing identity info on *gin.Context.
const (
	CtxUserID   = "auth.user_id"
	CtxUsername = "auth.username"
	CtxRole     = "auth.role"
	CtxClaims   = "auth.claims"
)

// Middleware returns a gin middleware that enforces JWT auth. If the service
// is disabled (no JWT secret), the middleware is a no-op so development
// environments stay in open-access mode.
//
// Paths supplied in Config.PublicPaths bypass the check. On success the
// middleware stores the user ID, username, role and claims in the gin.Context
// for downstream handlers.
func (s *TokenService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.Enabled() {
			c.Next()
			return
		}
		if s.IsPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
			})
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be 'Bearer <token>'",
			})
			return
		}
		claims, err := s.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}
		if claims.Kind != KindAccess {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "refresh token cannot be used as access token",
			})
			return
		}
		// Service tokens are not backed by a users row: IssueServiceToken
		// writes the caller-chosen subject straight into `sub` (e.g.
		// "service:risk"), which is not an int64. Requiring a numeric subject
		// here would 401 every legitimate service token — the role is the
		// authorisation, not a user id. No user id is fabricated: CtxUserID is
		// left unset so audit consumers see "no user" rather than user 0.
		if claims.Role == ServiceRole {
			c.Set(CtxUsername, claims.Username)
			c.Set(CtxRole, claims.Role)
			c.Set(CtxClaims, claims)
			c.Next()
			return
		}
		uid, err := claims.UserIDInt64()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "malformed token subject",
			})
			return
		}
		c.Set(CtxUserID, uid)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxClaims, claims)
		c.Next()
	}
}

// RequireRole is the auth-aware variant of the package-level RequireRole: it
// short-circuits to a no-op when the TokenService is not enabled.
//
// This matters because Middleware() is itself a no-op when Enabled()==false —
// it never sets CtxRole. Chaining the package-level RequireRole behind a
// disabled Middleware would therefore 401 every request, silently breaking
// the documented "auth disabled => open access" mode.
//
// Rule: whenever a route group is guarded by s.Middleware(), guard it with
// s.RequireRole(...) — not auth.RequireRole(...). The two decisions (is auth
// on? what role is needed?) must agree, and only the service knows the first
// answer.
//
// Place AFTER s.Middleware().
func (s *TokenService) RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.Enabled() {
			c.Next()
			return
		}
		RequireRole(allowed...)(c)
	}
}

// RequireRole returns a middleware that aborts with 403 if the caller's role
// is not in the allowed set. Place AFTER Middleware().
//
// Prefer TokenService.RequireRole for route groups guarded by
// TokenService.Middleware() — see its doc comment for why.
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := RoleFromContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authenticated user (RequireRole placed before Middleware?)",
			})
			return
		}
		for _, r := range allowed {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient role",
			"have":  role,
			"need":  allowed,
		})
	}
}

// UserFromContext returns the authenticated user ID and role, or false if the
// request is unauthenticated (or the middleware is disabled).
func UserFromContext(c *gin.Context) (int64, string, bool) {
	uidVal, uidOK := c.Get(CtxUserID)
	roleVal, roleOK := c.Get(CtxRole)
	if !uidOK || !roleOK {
		return 0, "", false
	}
	uid, ok1 := uidVal.(int64)
	role, ok2 := roleVal.(string)
	return uid, role, ok1 && ok2
}

// RoleFromContext returns just the role, or "" / false.
func RoleFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get(CtxRole)
	if !ok {
		return "", false
	}
	r, ok := v.(string)
	return r, ok
}

// UsernameFromContext returns the authenticated username, or "" / false.
func UsernameFromContext(c *gin.Context) (string, bool) {
	v, ok := c.Get(CtxUsername)
	if !ok {
		return "", false
	}
	u, ok := v.(string)
	return u, ok
}

// ClaimsFromContext returns the parsed claims, or nil / false.
func ClaimsFromContext(c *gin.Context) (*Claims, bool) {
	v, ok := c.Get(CtxClaims)
	if !ok {
		return nil, false
	}
	claims, ok := v.(*Claims)
	return claims, ok
}
