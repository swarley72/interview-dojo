package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const AuthClaimsKey contextKey = "AuthClaimsKey"

type AuthClaims struct {
	UserID      string
	IsSuperUser bool
}

func AuthMiddleware(jwtSecret []byte) func(next http.Handler) http.Handler {
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

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			isSuperUser, _ := claims["is_super_user"].(bool)

			ctx := context.WithValue(req.Context(), AuthClaimsKey, AuthClaims{
				UserID:      userID,
				IsSuperUser: isSuperUser,
			})
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !MustAuthClaims(r.Context()).IsSuperUser {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MustAuthClaims(ctx context.Context) AuthClaims {
	return ctx.Value(AuthClaimsKey).(AuthClaims)
}

func UserIDFromAuthClaims(ctx context.Context) string {
	return MustAuthClaims(ctx).UserID
}
