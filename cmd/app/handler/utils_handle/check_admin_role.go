package utils_handle

import (
	"net/http"
	"strings"
	"tf_ocg/utils"
)

func CheckAdminRole(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}

	tokenString := strings.Split(authHeader, " ")[1]

	role, err := utils.GetRoleFromToken(tokenString)
	if err != nil {
		return false
	}

	return role == "admin"
}
