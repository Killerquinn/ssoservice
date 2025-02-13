package postgresinit

import (
	"context"
	"fmt"
	"sso/internal/config"
	"time"

	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// settings of db init
const (
	maxConns          = 20
	connsMaxLifetime  = 120
	minConns          = 5
	connMaxIdleTime   = 20
	healthCheckPeriod = 30
)

var (
	ErrCannotConnectToDB = errors.New("cannot connect to database")
)

func NewPsqlDB(c *config.Config) (*pgxpool.Pool, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
	)

	cfg, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(maxConns)
	cfg.MinConns = minConns
	cfg.MaxConnLifetime = (connsMaxLifetime * time.Second)
	cfg.MaxConnIdleTime = (connMaxIdleTime * time.Second)
	cfg.HealthCheckPeriod = (healthCheckPeriod * time.Second)

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
