package pqfuncs

import (
	"context"
	"log/slog"
	"sso/internal/config"
	postgresinit "sso/internal/lib/postgresql"

	"github.com/jackc/pgx/v5/pgxpool"
)

var psqlDB *pgxpool.Pool

func Run(ctx context.Context, c *config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	var err error
	psqlDB, err = postgresinit.NewPsqlDB(c)
	if err != nil {
		log.Error("failed to run db", slog.Any("err", err))
		return nil, err
	}
	log.Info("db is running!")

	go func() {
		<-ctx.Done()
		Stop()
	}()
	return psqlDB, nil
}

func Stop() {
	if psqlDB != nil {
		psqlDB.Close()
	}
}
