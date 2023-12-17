package cart_handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	var requestBody struct {
		VariantID int32 `json:"variantId"`
		Quantity  int32 `json:"quantity"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	variantID := requestBody.VariantID
	quantity := requestBody.Quantity

	// Kiểm tra số lượng trước khi thêm vào giỏ hàng
	if valid, message := isQuantityValidAdd(variantID, quantity); !valid {
		res.JSON(w, http.StatusOK, map[string]string{"error": message})
		return
	}

	cart, err := dbms.AddToCart(userID, variantID, quantity)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to add item to cart"))
		return
	}

	res.JSON(w, http.StatusOK, cart)
}

func isQuantityValidAdd(variantID, quantity int32) (bool, string) {
	var variant models.Variant
	err := dbms.GetVariantById(&variant, variantID)
	if err != nil {
		return false, "Failed to retrieve variant information"
	}

	if quantity > variant.CountInStock {
		return false, "Invalid quantity. Quantity exceeds available stock."
	}

	return true, ""
}
