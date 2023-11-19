package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"tf_ocg/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		_, err := utils.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: "))
			return
		}
		next.ServeHTTP(w, r)
	}
}

func AuthAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := utils.VerifyToken(tokenString)
		mapClaims, ok := claims.(jwt.MapClaims)
		if ok {
			role := mapClaims["role"]
			if role != "admin" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("You don't have permission"))
				return
			}
		}
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token"))
			return
		}
		next.ServeHTTP(w, r)
	}
}
