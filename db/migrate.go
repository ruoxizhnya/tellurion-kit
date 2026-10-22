package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// Execer is the part of a pgx pool (or connection / transaction) the migration
// runner needs. Keeping it an interface makes the runner testable without a
// database and usable inside a caller-owned transaction.
type Execer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// Migrate executes statements in order and stops at the first failure. Blank
// statements are skipped.
//
// This is the mechanism only: the statements, their order and their
// idempotency (IF NOT EXISTS, guards, …) are the caller's DDL, and the caller
// decides whether to record progress. There is deliberately no migration table
// here — that shape is a data-model decision, not plumbing.
func Migrate(ctx context.Context, db Execer, statements []string) error {
	for i, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := db.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("db: migration statement %d failed: %w", i+1, err)
		}
	}
	return nil
}
