package auth_handle

import (
	"encoding/json"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode("Log out successfully")
}
