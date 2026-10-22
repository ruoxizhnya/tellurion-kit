package db

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Compile-time pin: a real pool satisfies the runner's port, so services can
// pass one directly (or a transaction, should they want the runner inside it).
var _ Execer = (*pgxpool.Pool)(nil)

type fakeExecer struct {
	calls []string
	errAt int // 1-based statement number that fails; 0 = never
	err   error
}

func (f *fakeExecer) Exec(_ context.Context, sql string, _ ...any) (pgconn.CommandTag, error) {
	f.calls = append(f.calls, sql)
	if f.errAt > 0 && len(f.calls) == f.errAt {
		return pgconn.CommandTag{}, f.err
	}
	return pgconn.CommandTag{}, nil
}

func TestMigrate_ExecutesInOrder(t *testing.T) {
	t.Parallel()
	f := &fakeExecer{}
	require.NoError(t, Migrate(context.Background(), f, []string{"first", "second", "third"}))
	assert.Equal(t, []string{"first", "second", "third"}, f.calls)
}

func TestMigrate_SkipsBlankStatements(t *testing.T) {
	t.Parallel()
	f := &fakeExecer{}
	require.NoError(t, Migrate(context.Background(), f, []string{"", "  \n\t", "real", ""}))
	assert.Equal(t, []string{"real"}, f.calls)
}

func TestMigrate_EmptySequenceIsNoop(t *testing.T) {
	t.Parallel()
	f := &fakeExecer{}
	require.NoError(t, Migrate(context.Background(), f, nil))
	assert.Empty(t, f.calls)
}

// TestMigrate_StopsAtFirstFailure pins both halves of the failure contract:
// statements after the failing one are not executed, and the error names the
// 1-based position so an operator can find the statement in the caller's
// sequence.
func TestMigrate_StopsAtFirstFailure(t *testing.T) {
	t.Parallel()
	boom := errors.New("syntax error at or near \"BROKEN\"")
	f := &fakeExecer{errAt: 2, err: boom}

	err := Migrate(context.Background(), f, []string{"a", "BROKEN", "c"})
	require.Error(t, err)
	assert.ErrorIs(t, err, boom, "the driver error must stay inspectable")
	assert.Contains(t, err.Error(), "statement 2")
	assert.Equal(t, []string{"a", "BROKEN"}, f.calls, "the runner must not continue past a failure")
}

// TestMigrate_NoMigrationTable asserts the scope decision mechanically: the
// runner issues exactly the caller's statements and nothing else — no
// bookkeeping INSERT/UPDATE is smuggled in.
func TestMigrate_NoMigrationTable(t *testing.T) {
	t.Parallel()
	f := &fakeExecer{}
	require.NoError(t, Migrate(context.Background(), f, []string{"CREATE TABLE IF NOT EXISTS t (id int)"}))
	assert.Len(t, f.calls, 1, "the mechanism must not add its own bookkeeping statements")
}
