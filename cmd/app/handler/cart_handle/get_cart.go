// Trong package cart_handle hoặc một file tương tự
package cart_handle

import (
	"errors"
	"fmt"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/cart_dto/response"
	"tf_ocg/cmd/app/handler/utils_handle"
	database "tf_ocg/pkg/database_manager"
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
	var totalQuantity int32
	var totalPrice float64

	discountCode := r.URL.Query().Get("discountCode")
	var discountAmount float64
	if discountCode != "" {
		discountAmount, err = dbms.ApplyDiscountForOrder(database.Db, cartItems, discountCode)
		if err != nil {
			res.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	for _, cartItem := range cartItems {
		var product models.Product // Tạo biến mới cho mỗi sản phẩm

		err := dbms.GetProductById(&product, cartItem.ProductID)
		if err != nil {
			fmt.Println("Error getting product:", err)
			res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin sản phẩm thất bại"))
			return
		}

		cartResponse := response.CartResponse{
			CartID:        cartItem.CartID,
			UserID:        cartItem.UserID,
			ProductID:     cartItem.ProductID,
			Quantity:      cartItem.Quantity,
			TotalPrice:    cartItem.TotalPrice,
			ProductDetail: product, // Sử dụng biến product mới tạo
		}

		cartResponses = append(cartResponses, cartResponse)
		totalQuantity += cartItem.Quantity
		totalPrice += cartItem.TotalPrice
	}

	totalProducts := len(cartResponses)

	summary := map[string]interface{}{
		"totalProducts": totalProducts,
		"totalQuantity": totalQuantity,
		"totalPrice":    totalPrice,
	}

	if discountCode != "" {
		totalPrice -= discountAmount
		summary["totalPrice"] = totalPrice
		summary["discountAmount"] = discountAmount
	}

	result := map[string]interface{}{
		"cartItems": cartResponses,
		"summary":   summary,
	}

	res.JSON(w, http.StatusOK, result)
}
