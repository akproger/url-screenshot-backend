package middleware

import (
	"net/http"
	"strings"

	"github.com/akproger/url-screenshot-backend/handlers"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

// AuthMiddleware проверяет JWT-токен в заголовке Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Forbidden: no token provided", http.StatusForbidden)
			return
		}

		// Убираем префикс "Bearer " для получения токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Парсим токен
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Forbidden: invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
