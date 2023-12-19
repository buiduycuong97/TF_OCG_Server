package router

import (
	"github.com/gorilla/mux"
	optionValue "tf_ocg/cmd/app/handler/option_value_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupOptionValueRoutes(r *mux.Router) {
	//authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", authAdminMiddleware(optionValue.CreateOptionValue)).Methods("POST")
	r.HandleFunc("/{id}", authAdminMiddleware(optionValue.DeleteOptionValue)).Methods("DELETE")
}
