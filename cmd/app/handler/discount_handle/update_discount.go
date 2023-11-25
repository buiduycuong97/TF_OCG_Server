package discount_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 32)
	discountID := int32(did)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var updatedDiscount models.Discount
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &updatedDiscount)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = dbms.UpdateDiscount(&updatedDiscount, discountID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, updatedDiscount)
}
