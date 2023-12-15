package router

import (
	"github.com/gorilla/mux"
	category "tf_ocg/cmd/app/handler/category_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupCategoryRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("", authAdminMiddleware(category.CreateCategory)).Methods("POST")
	r.HandleFunc("", category.GetListCategories).Methods("GET")
	r.HandleFunc("/{id}", category.GetCategory).Methods("GET")
	r.HandleFunc("/{id}", authAdminMiddleware(category.UpdateCategory)).Methods("PUT")
	r.HandleFunc("/{id}", authAdminMiddleware(category.DeleteCategory)).Methods("DELETE")
	r.HandleFunc("/search", authMiddleware(category.SearchCategories)).Methods("GET")
}
