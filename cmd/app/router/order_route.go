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
	r.HandleFunc("/api/orders", order_handle.CreateOrderHandler).Methods("POST")
	r.HandleFunc("/api/orders/{orderID}/capture", order_handle.CaptureOrderHandler).Methods("POST")
	r.HandleFunc("/complete", order_handle.CompleteOrderHandler).Methods("POST")
	r.HandleFunc("/request-cancel", authMiddleware(order_handle.RequestCancelOrderHandler)).Methods("POST")
	r.HandleFunc("/get-pending-orders", authMiddleware(order_handle.ViewPendingOrdersHandler)).Methods("GET")
	r.HandleFunc("/get-order-being-delivered-orders", authMiddleware(order_handle.ViewOrderBeingDeliveredHandler)).Methods("GET")
	r.HandleFunc("/get-complete-the-order-orders", authMiddleware(order_handle.ViewCompleteTheOrderHandler)).Methods("GET")
	r.HandleFunc("/get-cancelled-orders", authMiddleware(order_handle.ViewCancelledOrdersHandler)).Methods("GET")
	r.HandleFunc("", authAdminMiddleware(order_handle.GetAllOrder)).Methods("GET")
}
