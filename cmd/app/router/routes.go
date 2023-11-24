package router

import (
	"github.com/gorilla/mux"
)

func InitializeRoutes(server *mux.Router) {
	userRouter := server.PathPrefix("/users").Subrouter()
	productRouter := server.PathPrefix("/product").Subrouter()
	authRouter := server.PathPrefix("/auth").Subrouter()
	categoryRouter := server.PathPrefix("/category").Subrouter()
	cartRouter := server.PathPrefix("/cart").Subrouter()
	orderRouter := server.PathPrefix("/order").Subrouter()
	SetupUserRoutes(userRouter)
	SetupAuthRoutes(authRouter)
	SetupProductRoutes(productRouter)
	SetupCategoryRoutes(categoryRouter)
	SetupCartRoutes(cartRouter)
	SetupOrderRoutes(orderRouter)
}
