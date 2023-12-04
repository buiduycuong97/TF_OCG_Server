package router

import (
	"github.com/gorilla/mux"
	"tf_ocg/cmd/app/middleware"
)
import user "tf_ocg/cmd/app/handler/user_handle"

func SetupUserRoutes(r *mux.Router) {
	// Middleware cho các API cần xác thực
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", user.CreateUser).Methods("POST")
	r.HandleFunc("", authAdminMiddleware(user.GetUsers)).Methods("GET")
	r.HandleFunc("/{id}", authMiddleware(user.GetUser)).Methods("GET")
	r.HandleFunc("/{id}", user.UpdateUser).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(user.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/filter/search-user", authAdminMiddleware(user.SearchUsers)).Methods("GET")
	r.HandleFunc("/change-password", authMiddleware(user.ChangePassword)).Methods("POST")
}
