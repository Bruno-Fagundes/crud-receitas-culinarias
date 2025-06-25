package middleware

import (
	"context"
	"log" // Adicione esta linha
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Tipo de chave personalizado
type contextKey string

const userKey contextKey = "user"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("JWTMiddleware: Recebendo requisição") // Log 1
		authHeader := r.Header.Get("Authorization")
		log.Printf("JWTMiddleware: Authorization Header: %s\n", authHeader) // Log 2

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("JWTMiddleware: Token ausente ou formato inválido") // Log 3
			http.Error(w, "Token ausente ou inválido", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("JWTMiddleware: Token String: %s\n", tokenString) // Log 4

		claims := &jwt.MapClaims{}

		jwtSecret := os.Getenv("JWT_SECRET")
		log.Printf("JWTMiddleware: JWT_SECRET do .env: %s (length: %d)\n", jwtSecret, len(jwtSecret)) // Log 5

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil // Use a variável local jwtSecret
		})

		if err != nil || !token.Valid {
			log.Printf("JWTMiddleware: Erro de validação do token: %v, Token válido: %t\n", err, token.Valid) // Log 6
			http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		log.Println("JWTMiddleware: Token validado com sucesso!") // Log 7
		ctx := context.WithValue(r.Context(), userKey, (*claims)["username"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
