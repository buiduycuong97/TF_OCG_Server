package router

import (
	"github.com/gorilla/mux"
	discount "tf_ocg/cmd/app/handler/discount_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupDiscountRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", authAdminMiddleware(discount.CreateDiscountHandler)).Methods("POST")
	r.HandleFunc("/get-discount/get-by-code", discount.GetDiscountByDiscountCode).Methods("GET")
	r.HandleFunc("", authMiddleware(discount.GetAllDiscounts)).Methods("GET")
	r.HandleFunc("/{id}", authMiddleware(discount.GetDiscountByID)).Methods("GET")
	r.HandleFunc("/{id}", authAdminMiddleware(discount.UpdateDiscount)).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(discount.DeleteDiscount)).Methods("DELETE")
}
