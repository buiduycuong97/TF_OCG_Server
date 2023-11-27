package router

import (
	"github.com/gorilla/mux"
	order_handle "tf_ocg/cmd/app/handler/order_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupOrderRoutes(r *mux.Router) {
	authMiddleware := middleware.AuthMiddleware
	authAdminMiddleware := middleware.AuthAdmin
	r.HandleFunc("/checkout", authMiddleware(order_handle.CheckoutHandler)).Methods("POST")
	r.HandleFunc("/complete", authMiddleware(order_handle.CompleteOrderHandler)).Methods("POST")
	r.HandleFunc("/request-cancel", authMiddleware(order_handle.RequestCancelOrderHandler)).Methods("POST")
	r.HandleFunc("/get-pending-orders", authMiddleware(order_handle.ViewPendingOrdersHandler)).Methods("GET")
	r.HandleFunc("/get-order-being-delivered-orders", authMiddleware(order_handle.ViewOrderBeingDeliveredHandler)).Methods("GET")
	r.HandleFunc("/get-complete-the-order-orders", authMiddleware(order_handle.ViewCompleteTheOrderHandler)).Methods("GET")
	r.HandleFunc("/get-request-to-cancel-order-orders", authMiddleware(order_handle.ViewRequestToCancelOrderHandler)).Methods("GET")
	r.HandleFunc("/get-cancelled-orders", authMiddleware(order_handle.ViewCancelledOrdersHandler)).Methods("GET")
	r.HandleFunc("/accept-order", authAdminMiddleware(order_handle.AcceptOrderHandler)).Methods("POST")
	r.HandleFunc("/accept-cancel-request", authAdminMiddleware(order_handle.AdminAcceptCancelRequestHandler)).Methods("POST")
	r.HandleFunc("/decline-cancel-request", authAdminMiddleware(order_handle.AdminDeclineCancelRequestHandler)).Methods("POST")
	r.HandleFunc("/decline-order", authAdminMiddleware(order_handle.AdminDeclineOrderHandler)).Methods("POST")
}
