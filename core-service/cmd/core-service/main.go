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

	"buf.build/go/protovalidate"
	protovalidate_interceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/ilyakaznacheev/cleanenv"
	corepb "github.com/swarley72/interview-dojo/proto/core"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swarley72/interview-dojo/core-service/internal/repository"
	"github.com/swarley72/interview-dojo/core-service/internal/server"
	"github.com/swarley72/interview-dojo/core-service/internal/service"
	"google.golang.org/grpc"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	GRPCPort    int    `env:"GRPC_PORT" env-required:"true"`
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

	tagRepo := repository.NewPostgresTagRepository(pool)
	questionRepo := repository.NewPostgresQuestionRepository(pool)
	userProgressRepo := repository.NewPostgresUserProgressRepository(pool)

	tagSvc := service.NewTagService(tagRepo)
	questionSvc := service.NewQuestionService(questionRepo)
	userProgressSvc := service.NewUserProgressService(questionRepo, userProgressRepo)

	grpcServer := server.NewGRPCServer(questionSvc, tagSvc, userProgressSvc)

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed create validator %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(protovalidate_interceptor.UnaryServerInterceptor(validator)),
	)
	corepb.RegisterCoreServiceServer(s, grpcServer)

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
