package handlers

import (
	"encoding/json"
	"log" // Adicione para logs de depuração
	"net/http"
	"os"
	"time"

	"github.com/Bruno-Fagundes/crud-receitas-culinarias/models"
	"github.com/golang-jwt/jwt/v5"
)

// REMOVA OU COMENTE ESTA LINHA GLOBAL:
// var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Erro no JSON", http.StatusBadRequest)
		return
	}

	if creds.Username != "bruno" || creds.Password != "senha123" {
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// LEIA A CHAVE SECRETA AQUI DENTRO DA FUNÇÃO
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("ERRO: JWT_SECRET não está definida no ambiente para o handler de login!")
		http.Error(w, "Erro interno do servidor: Chave secreta não configurada", http.StatusInternalServerError)
		return
	}
	// log.Printf("LoginHandler: JWT_SECRET para assinatura: %s (length: %d)\n", jwtSecret, len(jwtSecret)) // Opcional para depuração

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret)) // Use a variável local jwtSecret
	if err != nil {
		log.Printf("Erro ao gerar token: %v\n", err) // Log de erro
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
