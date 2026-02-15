package middleware

import (
	"context"
	"net/http"
	"strings"

	userpb "github.com/swarley72/interview-dojo/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(userService userpb.UserServiceClient) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "no authorization header provided", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			res, err := userService.ValidateToken(req.Context(), &userpb.ValidateTokenRequest{AccessToken: token})
			if err != nil {
				s, _ := status.FromError(err)
				if s.Code() == codes.Unauthenticated {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(req.Context(), UserIDKey, res.UserId)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
