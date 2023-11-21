package router

import (
	"github.com/gorilla/mux"
)

func InitializeRoutes(server *mux.Router) {
	userRouter := server.PathPrefix("/users").Subrouter()
	productRouter := server.PathPrefix("/product").Subrouter()
	authRouter := server.PathPrefix("/auth").Subrouter()
	SetupUserRoutes(userRouter)
	SetupAuthRoutes(authRouter)
	SetupProductRoutes(productRouter)
}
