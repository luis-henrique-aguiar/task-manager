package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			fmt.Println("DEBUG: Header mal formatado")
			app.unauthorizedResponse(w)
			return
		}

		tokenString := headerParts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			fmt.Println("DEBUG: Token inválido ou assinatura errada:", err)
			app.unauthorizedResponse(w)
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			fmt.Printf("DEBUG: Erro no cast do sub. Valor: %v Tipo: %T\n", claims["sub"], claims["sub"])
			app.unauthorizedResponse(w)
			return
		}
		userID := int64(userIDFloat)

		r = app.contextSetUser(r, userID)
		next.ServeHTTP(w, r)
	})
}

func (app *application) unauthorizedResponse(w http.ResponseWriter) {
	http.Error(w, "Token inválido ou ausente", http.StatusUnauthorized)
}