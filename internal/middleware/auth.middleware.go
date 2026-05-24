package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"missing authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")

		if len(splitToken) != 2 {
			http.Error(
				w,
				"invalid token format",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := splitToken[1]

		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				return []byte(
					os.Getenv("SESSION_SECRET"),
				), nil
			},
		)

		if err != nil || !token.Valid {

			http.Error(
				w,
				"invalid token",
				http.StatusUnauthorized,
			)

			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			claims,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

// Cara pakai
// protectedHandler := http.HandlerFunc(
// 	authHandler.Profile,
// )

// http.Handle(
// 	"/profile",
// 	middleware.AuthMiddleware(
// 		protectedHandler,
// 	),
// )
