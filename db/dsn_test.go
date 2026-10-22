package db

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDSN_URLPassthrough(t *testing.T) {
	t.Parallel()
	const dsn = "postgres://alice:pw@db.internal:6543/prod?sslmode=require"
	got, err := BuildDSN(Config{URL: dsn, User: "ignored", Password: "ignored"})
	require.NoError(t, err)
	assert.Equal(t, dsn, got)
}

func TestBuildDSN_PlaceholderRejected(t *testing.T) {
	t.Parallel()
	_, err := BuildDSN(Config{URL: "postgres://u:${DB_PASSWORD}@h:5432/db"})
	require.ErrorIs(t, err, ErrDSNPlaceholder)
}

func TestBuildDSN_RequiresUserAndPassword(t *testing.T) {
	t.Parallel()
	_, err := BuildDSN(Config{Password: "pw"})
	assert.ErrorIs(t, err, ErrEmptyDBUser)

	_, err = BuildDSN(Config{User: "alice"})
	assert.ErrorIs(t, err, ErrEmptyDBPassword)

	// A URL is exempt from the field checks — a DSN may legitimately carry
	// credentials some other way (or none) and the driver will decide.
	got, err := BuildDSN(Config{URL: "postgres:///db?host=/var/run/postgresql"})
	require.NoError(t, err)
	assert.Equal(t, "postgres:///db?host=/var/run/postgresql", got)
}

func TestBuildDSN_Defaults(t *testing.T) {
	t.Parallel()
	got, err := BuildDSN(Config{User: "alice", Password: "pw"})
	require.NoError(t, err)
	assert.Equal(t, "postgres://alice:pw@localhost:5432/tellurion?sslmode=disable", got)
}

func TestBuildDSN_RespectsEveryField(t *testing.T) {
	t.Parallel()
	got, err := BuildDSN(Config{
		Host:     "db.internal",
		Port:     6543,
		User:     "alice",
		Password: "pw",
		Name:     "production",
		SSLMode:  "require",
	})
	require.NoError(t, err)
	assert.Equal(t, "postgres://alice:pw@db.internal:6543/production?sslmode=require", got)
}

// TestBuildDSN_EscapesPassword is the regression guard for the bug the single
// implementation exists to prevent: a password containing userinfo
// delimiters must not split the DSN. Parsing the result back with pgx proves
// the password survives byte-for-byte.
func TestBuildDSN_EscapesPassword(t *testing.T) {
	t.Parallel()
	const password = "p@ss:w/rd?#"
	dsn, err := BuildDSN(Config{Host: "h", User: "u", Password: password, Name: "n"})
	require.NoError(t, err)
	assert.NotContains(t, dsn, password, "the raw password must be percent-escaped, not embedded")

	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, password, cfg.ConnConfig.Password)
	assert.Equal(t, "u", cfg.ConnConfig.User)
	assert.Equal(t, "h", cfg.ConnConfig.Host)
	assert.Equal(t, uint16(5432), cfg.ConnConfig.Port)
	assert.Equal(t, "n", cfg.ConnConfig.Database)
}

func TestBuildDSN_IPv6Host(t *testing.T) {
	t.Parallel()
	dsn, err := BuildDSN(Config{Host: "::1", Port: 6000, User: "u", Password: "p"})
	require.NoError(t, err)

	cfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, "::1", cfg.ConnConfig.Host)
	assert.Equal(t, uint16(6000), cfg.ConnConfig.Port)
}
