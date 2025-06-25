package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("Senha12345")

func GerarToken(usuarioID string) (string, error) {
	claims := jwt.MapClaims{
		"usuario_id": usuarioID,
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
