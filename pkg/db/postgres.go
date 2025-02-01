package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

const (
	_defaultMaxPoolSize  = 2
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(connUrl string, opts ...Option) (*Postgres, error) {

	slog.With("db.New")
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	connConfig, err := pgxpool.ParseConfig(connUrl)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}
	connConfig.MaxConns = int32(pg.maxPoolSize)
	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		slog.InfoContext(ctx, "postgres - New - pgxpool.AfterConnect: connected!")
		return nil
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), connConfig)
		if err == nil {
			break
		}

		slog.Warn("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		slog.Error(err.Error())
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}

	if err = pg.Pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *Postgres) Close() error {
	if p.Pool != nil {
		p.Pool.Close()
	}

	return nil
}
