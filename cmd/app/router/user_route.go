package router

import (
	"github.com/gorilla/mux"
)
import user "tf_ocg/cmd/app/handler/user_handle"

func SetupUserRoutes(r *mux.Router) {
	r.HandleFunc("", user.CreateUser).Methods("POST")
	r.HandleFunc("", user.GetUsers).Methods("GET")
	r.HandleFunc("/{id}", user.GetUser).Methods("GET")
	r.HandleFunc("/{id}", user.UpdateUser).Methods("PUT")
	r.HandleFunc("/{id}", user.DeleteUser).Methods("DELETE")
}
