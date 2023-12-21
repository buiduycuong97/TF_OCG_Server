package discount_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
)

func DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 32)
	discountID := int32(did)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = dbms.DeleteDiscountIfExists(discountID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Discount deleted successfully"})
}
