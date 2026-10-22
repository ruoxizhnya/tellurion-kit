// Package db provides database plumbing: DSN assembly and validation,
// PostgreSQL / Redis connection setup, and a migration runner that executes a
// caller-supplied sequence of statements.
//
// It owns no schema. Every DDL statement, table name and column belongs to the
// service that owns that data (SCOPE.md §3 `db` row): this package never
// carries a CREATE TABLE and never records which migrations ran. Idempotency
// is the caller's contract — the runner simply executes statements in order
// and stops at the first failure.
package db

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Sentinel errors for BuildDSN. They are exported so callers and tests can
// tell the misconfigurations apart without matching on message text.
var (
	// ErrDSNPlaceholder 表示传入的 DSN 里还留着 `${...}` 占位符。
	ErrDSNPlaceholder = errors.New("database url contains an unresolved ${...} placeholder")
	// ErrEmptyDBUser 表示逐字段拼装时用户名是空的。
	ErrEmptyDBUser = errors.New("database user is empty")
	// ErrEmptyDBPassword 表示逐字段拼装时密码是空的。
	ErrEmptyDBPassword = errors.New("database password is empty")
)

// Config is the set of values needed to build a PostgreSQL DSN. It serves any
// service that accepts either a complete URL or discrete connection fields.
type Config struct {
	// URL is a complete DSN. When non-empty it is used as-is, and the fields
	// below are ignored.
	URL string

	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// BuildDSN resolves a Config into a PostgreSQL connection string.
//
// Contract:
//   - URL containing `${` → ErrDSNPlaceholder. There is no env expander: a
//     `${X}` in config is those characters verbatim, and feeding it to the
//     driver surfaces later as "password authentication failed", which reads
//     as a wrong password rather than a missing one.
//   - URL non-empty and placeholder-free → returned unchanged (the escape
//     hatch for embedders that already hold a full DSN).
//   - Otherwise assembled from the fields, with user and password required.
//     An empty password can never connect, so it is reported at startup
//     instead of at first query.
//
// Assembly goes through net/url rather than fmt.Sprintf: url.UserPassword
// percent-escapes by the userinfo rules, so a password containing `:` / `@` /
// `/` cannot split the DSN.
func BuildDSN(c Config) (string, error) {
	if strings.Contains(c.URL, "${") {
		return "", fmt.Errorf("%w: %s — there is no env expander, so it would be "+
			"used as the literal password; leave the url empty (it is then assembled "+
			"from the discrete fields) or supply a real DSN", ErrDSNPlaceholder, c.URL)
	}
	if c.URL != "" {
		return c.URL, nil
	}
	if c.User == "" {
		return "", fmt.Errorf("%w: set database.user (env: DATABASE_USER)", ErrEmptyDBUser)
	}
	if c.Password == "" {
		return "", fmt.Errorf("%w: set DATABASE_PASSWORD — a ${...} placeholder in "+
			"config does NOT work (there is no env expander)", ErrEmptyDBPassword)
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}
	port := c.Port
	if port == 0 {
		port = 5432
	}
	name := c.Name
	if name == "" {
		name = "tellurion"
	}
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(host, strconv.Itoa(port)),
		Path:     "/" + name,
		RawQuery: url.Values{"sslmode": {sslmode}}.Encode(),
	}
	return u.String(), nil
}
