package router

import (
	"github.com/gorilla/mux"
	optionValue "tf_ocg/cmd/app/handler/option_value_handle"
)

func SetupOptionValueRoutes(r *mux.Router) {
	//authMiddleware := middleware.AuthMiddleware
	//authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", optionValue.CreateOptionValue).Methods("POST")
	r.HandleFunc("/get-by/{id}", optionValue.GetOptionValueByOptionProductId).Methods("GET")
}
