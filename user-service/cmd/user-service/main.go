package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	protovalidate_interceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	userpb "github.com/swarley72/interview-dojo/proto/user"
	"github.com/swarley72/interview-dojo/user-service/internal/repository"
	"github.com/swarley72/interview-dojo/user-service/internal/server"
	"github.com/swarley72/interview-dojo/user-service/internal/service"
	"google.golang.org/grpc"
)

type Config struct {
	DatabaseURL     string        `env:"DATABASE_URL" env-required:"true"`
	GRPCPort        int           `env:"GRPC_PORT" env-required:"true"`
	JWTSecret       string        `env:"JWT_SECRET" env-required:"true"`
	AccessTokenExp  time.Duration `env:"ACCESS_TOKEN_EXP" env-default:"15m"`
	RefreshTokenExp time.Duration `env:"REFRESH_TOKEN_EXP" env-default:"720h"`
}

var cfg Config

func main() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("environment variable %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("failed connect to db - %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))

	if err != nil {
		log.Fatalf("failed to listen port %d, %v", cfg.GRPCPort, err)
	}

	repo := repository.NewPostgresUserRepository(pool)
	userService := service.NewUserService(
		repo, []byte(cfg.JWTSecret),
		cfg.AccessTokenExp,
		cfg.RefreshTokenExp,
	)
	grpcServer := server.NewGRPCServer(userService)

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed create validator %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(protovalidate_interceptor.UnaryServerInterceptor(validator)),
	)

	userpb.RegisterUserServiceServer(s, grpcServer)

	go func() {
		slog.Info("server started", "port", cfg.GRPCPort)
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("shutting down", "signal", sig.String())

	s.GracefulStop()
}
