package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgxutil"

	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

type txContextKeyType struct{}

var txContextKey txContextKeyType

type PgDoer interface {
	pgxutil.Execer
	pgxutil.Queryer
}

func (db *DB) Doer(ctx context.Context) PgDoer {
	tx := txFromContext(ctx)
	if tx != nil {
		return tx
	}
	return db.pool
}

func contextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func txFromContext(ctx context.Context) pgx.Tx {
	tx := ctx.Value(txContextKey)
	if tx != nil {
		return tx.(pgx.Tx)
	}
	return nil
}

func (db *DB) RunWithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	var err error

	rootTx := false

	tx := txFromContext(ctx)
	if tx == nil {
		tx, err = db.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to start transaction, %w", err)
		}

		rootTx = true
		ctx = contextWithTx(ctx, tx)
	}

	err = fn(ctx)
	if err != nil {
		if rootTx {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.WithError(log.FromContext(ctx), err).Error("failed to rollback transaction")
			}
		}
		return err
	}

	if rootTx {
		if err = tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transation, %w", err)
		}
	}

	return nil
}
