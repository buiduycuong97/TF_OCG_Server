package router

import (
	"github.com/gorilla/mux"
	review "tf_ocg/cmd/app/handler/review_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupReviewRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	r.HandleFunc("/add-review", authMiddleware(review.AddReviewHandler)).Methods("POST")
}
