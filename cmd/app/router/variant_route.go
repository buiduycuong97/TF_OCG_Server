package router

import (
	"github.com/gorilla/mux"
	variant "tf_ocg/cmd/app/handler/variant_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupVariantRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	r.HandleFunc("/add-variant", authMiddleware(variant.CreateVariantHandler)).Methods("POST")
	r.HandleFunc("/get-variant-id", authMiddleware(variant.GetVariantIdByOption)).Methods("GET")
	r.HandleFunc("/get-variant-by-order-id", variant.GetListVariantByOrderId).Methods("GET")
	r.HandleFunc("/update-product-quantity/{id}", authMiddleware(variant.UpdateVariantQuantityHandler)).Methods("PUT")
}
