package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/config"
	jaegerT "sso/internal/lib/jaeger"
	"sso/internal/lib/postgresql/pqfuncs"

	"github.com/opentracing/opentracing-go"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	//создание логгера

	//инициализация приложения(app)

	//запуск gRPC-сервера приложения

	//config init

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("cfg: ", cfg))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//pooling start
	psqlDB, err := pqfuncs.Run(ctx, cfg, log)
	if err != nil || psqlDB == nil {
		log.Error("Postgres init error", slog.Any("err", err))
	}

	//opentracing, jaeger init
	tracer, closer, err := jaegerT.InitJaeger(cfg)
	if err != nil {
		log.Error("Init jaeger error", slog.Any("err", err))
	}
	log.Info("jaeger connected")

	opentracing.SetGlobalTracer(tracer)

	defer closer.Close()
	//starting application
	application := app.New(log, cfg.GRPC.Port, &cfg.Postgres, cfg.TokenTTL, psqlDB)

	go application.GRPCServer.MustRun()

	//signal monitoring
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	siginf := <-stop

	log.Info("stopping aplication", slog.String("last signal", siginf.String()))

	application.GRPCServer.Stop()
	pqfuncs.Stop()

	log.Info("application will stop after manage last orders before signal")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}
	return log
}
