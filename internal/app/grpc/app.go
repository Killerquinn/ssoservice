package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	authgrpc "sso/internal/grpc/auth"
	"sso/internal/interceptors"
	jwtlib "sso/internal/lib/jwt"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, authService authgrpc.AuthS) *App {

	jwt_new, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Warn("cannot fint jwt secret variable")

		return nil
	}

	authSvc, err := jwtlib.NewService(jwt_new)
	if err != nil {
		log.Warn("cannot register new service")

		return nil
	}

	interceptor, err := interceptors.NewAuthInterceptor(authSvc)
	if err != nil {
		log.Warn("cannot define interceptors.NewAuthInterceptor")

		return nil
	}
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryAuthInterceptor))

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Runs app without error handler to user
//
// Because there is no sence to start app if there will be an error
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Runs app
func (a *App) Run() error {
	const op = "grpcapp.Run" //op means operation, constant for more gentle method of logger

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("grpc successfully running!", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("grpc server stops serve")

	a.gRPCServer.GracefulStop()
}
