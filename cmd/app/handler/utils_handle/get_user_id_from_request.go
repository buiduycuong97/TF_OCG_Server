package utils_handle

import (
	"errors"
	"net/http"
	"strings"
	"tf_ocg/utils"
)

func GetUserIDFromRequest(r *http.Request) (int32, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return 0, errors.New("Missing Authorization header")
	}

	tokenString := strings.Split(authHeader, " ")[1]

	userID, err := utils.GetUserFromToken(tokenString)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
