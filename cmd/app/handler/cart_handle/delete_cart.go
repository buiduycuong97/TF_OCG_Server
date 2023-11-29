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

func RemoveCartItemHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid product ID"))
		return
	}

	err = dbms.RemoveCartItem(userID, productID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to remove item from cart"))
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Item removed from cart successfully"})
}
