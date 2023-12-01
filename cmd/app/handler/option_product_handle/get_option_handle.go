package option_product_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func GetOptionProductByProductId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])

	optionProductList, err := dbms.GetOptionProductByProductId(int32(productId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.JSON(w, http.StatusOK, optionProductList)
}

func GetAllOptionProduct(w http.ResponseWriter, r *http.Request) {
	optionProductList, err := dbms.GetAllOptionProduct()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.JSON(w, http.StatusOK, optionProductList)
}
