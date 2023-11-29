package router

import (
	"github.com/gorilla/mux"
	"tf_ocg/cmd/app/handler/province_handle"
)

func SetupProvinceRoutes(r *mux.Router) {
	r.HandleFunc("/get-province-by-name", province_handle.GetProvinceByName).Methods("GET")
	r.HandleFunc("/get-province-by-id", province_handle.GetProvinceById).Methods("GET")
	r.HandleFunc("", province_handle.GetAllProvince).Methods("GET")
}
