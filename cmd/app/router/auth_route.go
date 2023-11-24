package router

import (
	"github.com/gorilla/mux"
	auth "tf_ocg/cmd/app/handler/auth_handle"
)

func SetupAuthRoutes(r *mux.Router) {
	r.HandleFunc("", auth.HandleHome)
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/login-admin", auth.LoginAdmin).Methods("POST")
	r.HandleFunc("/refresh-token", auth.RefreshToken).Methods("GET")
	r.HandleFunc("/logout", auth.Logout)
	r.HandleFunc("/login-google", auth.HandleLogin)
	r.HandleFunc("/callback-google", auth.HandleCallback)
	r.HandleFunc("/forget-password", auth.HandleForgetPassword).Methods("POST")
	r.HandleFunc("/reset-password", auth.HandleResetPassword).Methods("POST")
	//r.HandleFunc("/login-facebook", auth.FacebookProvider.HandleLogin)
	//r.HandleFunc("/callback-facebook", auth.FacebookProvider.HandleCallback)
}
