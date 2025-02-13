package app

import (
	"log/slog"
	"os"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	jwtlib "sso/internal/lib/jwt"
	"sso/internal/services/authsvc"
	userrepository "sso/internal/storage/repository/auth_repo"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, postgres *config.PostgresConfig, tokenTTL time.Duration, db *pgxpool.Pool) *App {
	//init storageg
	dsn := os.Getenv("POSTGRES_URL")
	storage, err := userrepository.New(dsn)
	if err != nil {
		panic(err)
	}
	secret := os.Getenv("SECRET_JWT")
	tokengen, err := jwtlib.NewService(secret)
	//init auth service(auth)

	authService := authsvc.New(log, storage, storage, storage, tokengen, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
