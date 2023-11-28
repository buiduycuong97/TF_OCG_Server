package router

import (
	"github.com/gorilla/mux"
	"tf_ocg/cmd/app/handler/province_handle"
)

func SetupProvinceRoutes(r *mux.Router) {

	r.HandleFunc("", province_handle.GetAllProvince).Methods("GET")
}
