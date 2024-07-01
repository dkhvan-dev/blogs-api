package middlewares

import (
	"blogs-api/internal/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

func AuthMiddleware(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := verifyToken(r)
		if err != nil || !token.Valid {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		currentUserRole := claims["role"]
		if role != "" && currentUserRole != role {
			http.Error(w, `{"error": "Доступ запрещен."}`, http.StatusForbidden)
			return
		}

		context.Set(r, "login", claims["iss"])
		context.Set(r, "role", claims["role"])
		next.ServeHTTP(w, r)
	})
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	bearerToken := r.Header.Get("Authorization")
	tokenString := ""
	if strings.HasPrefix(bearerToken, "Bearer ") {
		tokenString = strings.TrimPrefix(bearerToken, "Bearer ")
	}
	secretKey := []byte(utils.GetEnv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
