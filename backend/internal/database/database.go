package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// DB wraps sqlx.DB providing convenience methods with context timeouts.
type DB struct {
	*sqlx.DB
}

const defaultQueryTimeout = 10 * time.Second

// QueryContext executes a query with a default timeout context derived from the parent.
func (db *DB) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultQueryTimeout)
}
