package router

import (
	"github.com/gorilla/mux"
	cart "tf_ocg/cmd/app/handler/cart_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupCartRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	r.HandleFunc("/add-to-cart", authMiddleware(cart.AddToCartHandler)).Methods("POST")
	r.HandleFunc("/view-cart", authMiddleware(cart.ViewCartHandler)).Methods("GET")
	r.HandleFunc("/remove-cart-item/{variantId}", authMiddleware(cart.RemoveCartItemHandler)).Methods("DELETE")
	r.HandleFunc("/update-cart-item/{variantId}", authMiddleware(cart.UpdateCartItemHandler)).Methods("PUT")
}
