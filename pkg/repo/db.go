package repo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/gavrilaf/wardrobe/pkg/utils"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

const (
	connectRetries = 10
	pingTimeout    = 3 * time.Second
	logStatPeriod  = 5 * 60 * time.Second // every 5 minutes
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrUniqueConstraint = errors.New("duplicate key value violates unique constraint")
)

type DB struct {
	pool *pgxpool.Pool
	quit chan struct{}
}

type TxRunner interface {
	RunWithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewDB(ctx context.Context, connString string, maxConn int32) (*DB, error) {
	// Configure the driver to connect to the database
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %v", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.ConnConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
	}

	config.MaxConns = maxConn
	config.MinConns = 2

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   log.PgxLogAdapter{},
		LogLevel: tracelog.LogLevelInfo,
	}

	var pool *pgxpool.Pool
	operation := func() error {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			return backoff.Permanent(err)
		}

		// pool created as lazy, so we must ping db to acquire the real connection
		err = pool.Ping(ctx)
		if err != nil {
			log.WithError(log.FromContext(ctx), err).Info("failed to ping database")
		}

		return err
	}

	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), connectRetries))
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool, %w", err)
	}

	quit := make(chan struct{}, 1)

	// log database statistic every 5 minutes
	utils.FnTicker(logStatPeriod, quit, func() {
		stat := pool.Stat()
		log.FromContext(ctx).Infof("db stat(total=%d, acquired=%d, construction=%d, idle=%d)",
			stat.TotalConns(), stat.AcquiredConns(), stat.ConstructingConns(), stat.IdleConns())
	})

	return &DB{pool: pool, quit: quit}, nil
}

func (db *DB) Migrate(ctx context.Context, path string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	sourceURL := "file:///" + filepath.Join(dir, path)
	log.FromContext(ctx).Infof("Looking for migration scripts in: %s", sourceURL)

	connStr := db.pool.Config().ConnString()

	m, err := migrate.New(sourceURL, connStr)
	if err != nil {
		return err
	}

	if err = m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	return db.pool.Ping(ctx)
}

func (db *DB) Close(ctx context.Context) {
	db.quit <- struct{}{}
	db.pool.Close()
	log.FromContext(ctx).Info("database closed")
}

func IsNoRowsFound(err error) bool {
	return err.Error() == pgx.ErrNoRows.Error()
}
