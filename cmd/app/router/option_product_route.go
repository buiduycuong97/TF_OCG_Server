package router

import (
	"github.com/gorilla/mux"
	optionProduct "tf_ocg/cmd/app/handler/option_product_handle"
)

func SetupOptionProductRoutes(r *mux.Router) {
	//authMiddleware := middleware.AuthMiddleware
	//authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", optionProduct.CreateOptionProductHandler).Methods("POST")
	r.HandleFunc("", optionProduct.GetAllOptionProduct).Methods("GET")
	r.HandleFunc("/get-by/{id}", optionProduct.GetOptionProductByProductId).Methods("GET")
}
