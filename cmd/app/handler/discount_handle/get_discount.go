package discount_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func GetDiscountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 32)
	discountID := int32(did)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var discount models.Discount
	err = dbms.GetDiscountByID(&discount, discountID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, discount)
}

func GetAllDiscounts(w http.ResponseWriter, r *http.Request) {
	discounts, err := dbms.GetAllDiscounts()
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, discounts)
}