package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/swarley72/interview-dojo/api-gateway/internal/handler"
	"github.com/swarley72/interview-dojo/api-gateway/internal/middleware"
	userpb "github.com/swarley72/interview-dojo/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	UserServiceAddr string `env:"USER_SERVICE_ADDR" env-required:"true"`
	JWTSecret       string `env:"JWT_SECRET" env-required:"true"`
	Port            int    `env:"HTTP_PORT" env-required:"true"`
}

var cfg Config

func main() {
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("environment variable %v", err)
	}

	userServiceConn, err := grpc.NewClient(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cannot connect to userServce %v", err)
	}

	userServiceClient := userpb.NewUserServiceClient(userServiceConn)

	h := handler.NewHandler(userServiceClient, validator.New())

	r := chi.NewRouter()
	r.Post("/api/register", h.Register)
	r.Post("/api/login", h.Login)
	r.Post("/api/refresh-token", h.RefreshToken)
	r.With(middleware.AuthMiddleware(userServiceClient)).Get("/api/profile", h.GetProfile)

	go func() {
		slog.Info("serve listening", "port", cfg.Port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)

		if err != nil {
			log.Fatalf("serve listening error %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("shutting down", "signal", sig.String())

}
