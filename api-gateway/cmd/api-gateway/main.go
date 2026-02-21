package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/swarley72/interview-dojo/api-gateway/internal/handler"
	"github.com/swarley72/interview-dojo/api-gateway/internal/middleware"
	corepb "github.com/swarley72/interview-dojo/proto/core"
	userpb "github.com/swarley72/interview-dojo/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	UserServiceAddr string `env:"USER_SERVICE_ADDR" env-required:"true"`
	CoreServiceAddr string `env:"CORE_SERVICE_ADDR" env-required:"true"`
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
	defer userServiceConn.Close()

	coreServiceConn, err := grpc.NewClient(
		cfg.CoreServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cannot connect to coreServce %v", err)
	}
	defer coreServiceConn.Close()

	userServiceClient := userpb.NewUserServiceClient(userServiceConn)
	coreServiceClient := corepb.NewCoreServiceClient(coreServiceConn)

	h := handler.NewHandler(userServiceClient, coreServiceClient, validator.New())

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		// Публичные
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh-token", h.RefreshToken)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Требуют авторизации
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret)))
			r.Get("/profile", h.GetProfile)
			r.Get("/tags", h.ListTags)
			r.Get("/questions", h.ListQuestions)
			r.Get("/questions/{id}", h.GetQuestion)
			r.Get("/anki/next", h.GetNextQuestion)
			r.Post("/anki/{question_id}/answer", h.RecordAnswer)

			// Admin
			r.Group(func(r chi.Router) {
				r.Use(middleware.AdminOnly)

				r.Post("/questions", h.CreateQuestion)
				r.Patch("/questions/{id}", h.UpdateQuestion)
				r.Delete("/questions/{id}", h.DeleteQuestion)

				r.Post("/tags", h.CreateTag)
				r.Delete("/tags/{id}", h.DeleteTag)
			})
		})
	})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		slog.Info("serve listening", "port", cfg.Port)
		err := srv.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			slog.Error("serve listening error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("shutting down", "signal", sig.String())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
	}

}
