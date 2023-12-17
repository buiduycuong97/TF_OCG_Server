package cart_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	vars := mux.Vars(r)
	variantID, err := strconv.Atoi(vars["variantId"])
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid variant ID"))
		return
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid quantity"))
		return
	}

	// Kiểm tra số lượng trước khi cập nhật giỏ hàng
	if !isQuantityValid(int32(variantID), int32(quantity)) {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid quantity"))
		return
	}

	err = dbms.UpdateCartItem(int(userID), variantID, quantity)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to update cart"))
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Cart updated successfully"})
}

func isQuantityValid(variantID, quantity int32) bool {
	var variant models.Variant
	err := dbms.GetVariantById(&variant, variantID)
	if err != nil {
		return false
	}

	return quantity <= variant.CountInStock
}
