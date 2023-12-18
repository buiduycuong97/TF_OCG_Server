package router

import (
	"github.com/gorilla/mux"
	variant "tf_ocg/cmd/app/handler/variant_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupVariantRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin
	r.HandleFunc("/add-variant", authAdminMiddleware(variant.CreateVariantHandler)).Methods("POST")
	r.HandleFunc("/get-variant-id", authMiddleware(variant.GetVariantIdByOption)).Methods("POST")
	r.HandleFunc("/get-variant-by-order-id", variant.GetListVariantByOrderId).Methods("GET")
	r.HandleFunc("/update-product-quantity/{id}", authMiddleware(variant.UpdateVariantQuantityHandler)).Methods("PUT")
	r.HandleFunc("/get-variant/{id}", variant.GetVariantById).Methods("GET")
	r.HandleFunc("/update-variant/{id}", variant.UpdateVariantByAdmin).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(variant.DeleteVariant)).Methods("DELETE")

}
