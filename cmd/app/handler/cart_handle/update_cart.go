package cart_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
)

func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid product ID"))
		return
	}
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid quantity"))
		return
	}

	err = dbms.UpdateCartItem(userID, productID, quantity)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to update cart"))
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Cart updated successfully"})
}
