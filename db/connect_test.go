package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect_BadDSN(t *testing.T) {
	t.Parallel()
	_, err := Connect(context.Background(), "postgres://u:p@h:not-a-port/db", PoolOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse connection string")
}

// TestConnect_CanceledContextFailsCleanly pins the fail-loud half of Connect:
// an unusable context is reported as a ping failure, not a nil pool.
func TestConnect_CanceledContextFailsCleanly(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	pool, err := Connect(ctx, "postgres://u:p@127.0.0.1:5432/db?sslmode=disable", PoolOptions{})
	require.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "failed to ping database")
}

func TestConnectRedis_BadURL(t *testing.T) {
	t.Parallel()
	_, err := ConnectRedis(context.Background(), "://not-a-url")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse redis url")
}

func TestConnectRedis_CanceledContextFailsCleanly(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	client, err := ConnectRedis(ctx, "redis://127.0.0.1:6379/0")
	require.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "failed to ping redis")
}

func TestDefaultPoolOptions(t *testing.T) {
	t.Parallel()
	opts := DefaultPoolOptions()
	assert.Equal(t, int32(20), opts.MaxConns)
	assert.Equal(t, int32(5), opts.MinConns)
	assert.Equal(t, time.Hour, opts.MaxConnLifetime)
	assert.Equal(t, 30*time.Minute, opts.MaxConnIdleTime)
}
