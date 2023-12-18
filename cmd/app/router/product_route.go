package router

import (
	"github.com/gorilla/mux"
	product "tf_ocg/cmd/app/handler/product_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupProductRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", authAdminMiddleware(product.CreateProduct)).Methods("POST")
	r.HandleFunc("", authAdminMiddleware(product.GetListProducts)).Methods("GET")
	r.HandleFunc("/{id}", authMiddleware(product.GetProduct)).Methods("GET")
	r.HandleFunc("/find-product/handle", product.GetProductByHandle).Methods("GET")
	r.HandleFunc("", authAdminMiddleware(product.UpdateProduct)).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(product.DeleteProduct)).Methods("DELETE")
	r.HandleFunc("/category/get-list", authAdminMiddleware(product.GetListProductByCategoryId)).Methods("GET")
	r.HandleFunc("/search/list", product.SearchProducts).Methods("GET")

}
