package grpc

import (
	grpcauthentication "authorization-service/internal/grpc/authentication"
	serviceauthentication "authorization-service/internal/service/authentication"
	"fmt"
	"log/slog"
	"net"

	pgstorage "authorization-service/internal/storage/postgres"

	authorizationservicev1 "github.com/GrishanyaaShustov/CloudStorage-Protos-Service/gen/go/authorization-service"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App holds gRPC server instance and its configuration.
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	gRPCPort   int
}

// New creates a new gRPC server app but does NOT start it.
// The caller is responsible for running and stopping the server.
func New(log *slog.Logger, gRPCPort int, pg *pgxpool.Pool) *App {
	gRPCServer := grpc.NewServer()

	// Enable reflection for grpcurl / Postman
	reflection.Register(gRPCServer)

	// Wire repository
	userRepo := pgstorage.NewUserRepository(log, pg)

	// Wire authentication service: business layer + transport layer.
	authenticationService := serviceauthentication.NewAuthService(log, userRepo)
	authenticationServer := grpcauthentication.NewServer(log, authenticationService)

	// Register gRPC handler for AuthenticationService.
	authorizationservicev1.RegisterAuthenticationServiceServer(gRPCServer, authenticationServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		gRPCPort:   gRPCPort,
	}
}

// MustRun starts gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "grpcApp.Run"

	l, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", a.gRPCPort))
	if err != nil {
		return fmt.Errorf("%d:%w", op, err)
	}

	a.log.Info("gRPC server started",
		slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gracefully stops gRPC server.
func (a *App) Stop() {
	const op = "grpcApp.Stop"
	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.gRPCPort))

	a.gRPCServer.GracefulStop()
}
