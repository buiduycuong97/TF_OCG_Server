package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func GenerateAccessToken(id int32) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}
	secretKey := os.Getenv("JWT_SECRET")
	signingKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 1).Unix(), // Token hết hạn sau 1h
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func GenerateRefreshToken(id int32) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}
	secretKey := os.Getenv("JWT_SECRET")
	signingKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}
	secretKey := os.Getenv("JWT_SECRET")
	signingKey := []byte(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func GenerateAccessTokenAdmin(role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Lỗi khi tải tệp .env")
	}
	secretKey := os.Getenv("JWT_SECRET")
	signingKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Minute * 1).Unix(), // Token hết hạn sau 1h
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
func GenerateRefreshTokenAdmin(role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	signingKey := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err

}
