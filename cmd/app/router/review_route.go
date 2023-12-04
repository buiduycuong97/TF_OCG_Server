package router

import (
	"github.com/gorilla/mux"
	review "tf_ocg/cmd/app/handler/review_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupReviewRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	//authAdminMiddleware := middleware.AuthAdmin
	r.HandleFunc("/add-review", authMiddleware(review.AddReviewHandler)).Methods("POST")
	r.HandleFunc("/get-review-by-product-id", review.GetListReviewByProductID).Methods("GET")
}
