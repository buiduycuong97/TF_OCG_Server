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
	reviewRouter := server.PathPrefix("/review").Subrouter()
	discountRouter := server.PathPrefix("/discount").Subrouter()
	variantRouter := server.PathPrefix("/variant").Subrouter()
	optionProductRouter := server.PathPrefix("/option-product").Subrouter()
	optionValueRouter := server.PathPrefix("/option-value").Subrouter()
	optionProvince := server.PathPrefix("/province").Subrouter()
	orderDetailRouter := server.PathPrefix("/order-detail").Subrouter()
	transactionRouter := server.PathPrefix("/transaction").Subrouter()

	SetupUserRoutes(userRouter)
	SetupAuthRoutes(authRouter)
	SetupProductRoutes(productRouter)
	SetupCategoryRoutes(categoryRouter)
	SetupCartRoutes(cartRouter)
	SetupOrderRoutes(orderRouter)
	SetupReviewRoutes(reviewRouter)
	SetupDiscountRoutes(discountRouter)
	SetupVariantRoutes(variantRouter)
	SetupOptionProductRoutes(optionProductRouter)
	SetupOptionValueRoutes(optionValueRouter)
	SetupProvinceRoutes(optionProvince)
	SetupOrderDetailRoutes(orderDetailRouter)
	SetupTransactionRoutes(transactionRouter)
}
