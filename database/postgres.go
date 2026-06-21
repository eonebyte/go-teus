package database

import (
	"fmt"
	"time"

	"github.com/eonebyte/go-teus/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("database: connect: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdletime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	return db, nil
}

// WithTransaction runs fn inside a database transaction.
// Commits on success, rolls back on error or panic.
func WithTransaction(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Querier is implemented by both *sqlx.DB and *sqlx.Tx,
// allowing repositories to accept either.
type Querier interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	BinderContext
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (interface {
		LastInsertId() (int64, error)
		RowsAffected() (int64, error)
	}, error)
}

// BinderContext is the sqlx binder interface subset.
type BinderContext interface {
	BinderContext()
}

// Timeout constants for database operations.
const (
	DefaultQueryTimeout = 5 * time.Second
	DefaultExecTimeout  = 10 * time.Second
)
