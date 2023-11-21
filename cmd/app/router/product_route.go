package router

import (
	"github.com/gorilla/mux"
	product "tf_ocg/cmd/app/handler/product_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupProductRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", authMiddleware(product.CreateProduct)).Methods("POST")
	r.HandleFunc("", authMiddleware(product.GetListProducts)).Methods("GET")
	r.HandleFunc("/{id}", authMiddleware(product.GetProduct)).Methods("GET")
	r.HandleFunc("/{id}", authAdminMiddleware(product.UpdateProduct)).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(product.DeleteProduct)).Methods("DELETE")
	r.HandleFunc("/category/{categoryID}", authMiddleware(product.GetListProductByCategoryId)).Methods("GET")
	r.HandleFunc("/search", authMiddleware(product.SearchProducts)).Methods("GET")
}
