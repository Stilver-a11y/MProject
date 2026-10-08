package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)
var JwtSecret = []byte("change-me-to-a-long-random-secret")

type ContextKey string

const UserIDKey ContextKey = "user_id"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("authorization")
		if authHeader == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return 
	}

	parts := strings.Split(authHeader," ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return 
	}

	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error){
		return JwtSecret, nil
	}) 
	if err != nil || !token.Valid{
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return 	
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userID := uint(userIDFloat)
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
		return 
	}

	http.Error(w, "Invalid token claims", http.StatusUnauthorized)
	}
}