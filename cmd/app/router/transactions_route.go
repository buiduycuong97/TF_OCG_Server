package router

import (
	"github.com/gorilla/mux"
	transaction "tf_ocg/cmd/app/handler/transaction_handle"
	"tf_ocg/cmd/app/middleware"
)

func SetupTransactionRoutes(r *mux.Router) {
	// Middleware cho các API cần xác thực
	authMiddleware := middleware.AuthMiddleware

	r.HandleFunc("/create-transaction", authMiddleware(transaction.CreateTransaction)).Methods("POST")

}
