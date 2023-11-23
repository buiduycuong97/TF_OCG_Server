package cart_handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
)

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	var requestBody struct {
		ProductID int32 `json:"productId"`
		Quantity  int32 `json:"quantity"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	ProductID := requestBody.ProductID
	Quantity := requestBody.Quantity

	cart, err := dbms.AddToCart(userID, ProductID, Quantity)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to add item to cart"))
		return
	}

	res.JSON(w, http.StatusOK, cart)
}
