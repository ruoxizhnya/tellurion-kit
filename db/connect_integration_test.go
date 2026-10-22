//go:build integration

// Real-database checks for the connection and migration mechanisms.
//
// Both tests are opt-in: set TEST_DATABASE_DSN / TEST_REDIS_URL to run them.
// The migration leg works inside a transaction that is always rolled back and
// only touches a TEMP table, so it leaves no residue in the target database.
package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect_RealPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN 未设置；跳过真库连接测试")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := Connect(ctx, dsn, DefaultPoolOptions())
	require.NoError(t, err)
	defer pool.Close()
	require.NoError(t, pool.Ping(ctx))
}

func TestConnectRedis_RealRedis(t *testing.T) {
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("TEST_REDIS_URL 未设置；跳过真 Redis 连接测试")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ConnectRedis(ctx, url)
	require.NoError(t, err)
	require.NoError(t, client.Close())
}

func TestMigrate_RealPostgresInRolledBackTx(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN 未设置；跳过真库迁移执行测试")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := Connect(ctx, dsn, PoolOptions{MaxConns: 1})
	require.NoError(t, err)
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	require.NoError(t, err)
	rolledBack := false
	defer func() {
		if !rolledBack {
			_ = tx.Rollback(ctx)
		}
	}()

	require.NoError(t, Migrate(ctx, tx, []string{
		"CREATE TEMP TABLE kit_db_migrate_probe (id INT PRIMARY KEY)",
		"INSERT INTO kit_db_migrate_probe (id) VALUES (1), (2)",
	}))
	var n int
	require.NoError(t, tx.QueryRow(ctx, "SELECT COUNT(*) FROM kit_db_migrate_probe").Scan(&n))
	assert.Equal(t, 2, n, "statements must have executed in order in the same session")

	err = Migrate(ctx, tx, []string{"THIS IS NOT SQL"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "statement 1")

	require.NoError(t, tx.Rollback(ctx))
	rolledBack = true
}
