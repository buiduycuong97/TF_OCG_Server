package cart_handle

import (
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/cart_dto/response"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func ViewCartHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Token không hợp lệ"))
		return
	}

	cartItems, err := dbms.GetCartByUserID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy giỏ hàng thất bại"))
		return
	}

	var cartResponses []response.CartResponse
	var product models.Product
	for _, cartItem := range cartItems {
		err := dbms.GetProductById(&product, cartItem.ProductID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin sản phẩm thất bại"))
			return
		}

		cartResponse := response.CartResponse{
			CartID:        cartItem.CartID,
			UserID:        cartItem.UserID,
			ProductID:     cartItem.ProductID,
			Quantity:      cartItem.Quantity,
			TotalPrice:    cartItem.TotalPrice,
			ProductDetail: product,
		}

		cartResponses = append(cartResponses, cartResponse)
	}

	res.JSON(w, http.StatusOK, cartResponses)
}
