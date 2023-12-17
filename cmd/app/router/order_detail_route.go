package router

import (
	"github.com/gorilla/mux"
	orderDetail "tf_ocg/cmd/app/handler/order_detail_handle"
)

func SetupOrderDetailRoutes(r *mux.Router) {
	//authMiddleware := middleware.AuthMiddleware
	//authAdminMiddleware := middleware.AuthAdmin

	r.HandleFunc("/get-list-order-detail", orderDetail.GetOrderInfoHandler).Methods("GET")
}
