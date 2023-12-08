// Trong package cart_handle hoặc một file tương tự
package cart_handle

import (
	"errors"
	"fmt"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/cart_dto/response"
	variantresponse "tf_ocg/cmd/app/dto/variant_dto/response"
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
	var totalQuantity int32
	var totalPrice float64

	discountCode := r.URL.Query().Get("discountCode")
	var discountAmount float64
	if discountCode != "" {
		discountAmount, err = dbms.ApplyDiscountForOrder(cartItems, discountCode)
		if err != nil {
			res.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	for _, cartItem := range cartItems {
		var variant models.Variant
		var product models.Product
		var option1 models.OptionValue
		var option2 models.OptionValue

		// Lấy thông tin biến thể
		err := dbms.GetVariantById(&variant, cartItem.VariantID)
		if err != nil {
			fmt.Println("Error getting variant:", err)
			res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin biến thể thất bại"))
			return
		}

		// Lấy thông tin sản phẩm từ ProductID trong biến thể
		err = dbms.GetProductById(&product, variant.ProductID)
		if err != nil {
			fmt.Println("Error getting product:", err)
			res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin sản phẩm thất bại"))
			return
		}

		if variant.OptionValue1 != 0 {
			err = dbms.GetOptionValueById(&option1, variant.OptionValue1)
			if err != nil {
				fmt.Println("Error getting option 1:", err)
				res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin option 1 thất bại"))
				return
			}
		}

		if variant.OptionValue2 != 0 {
			err = dbms.GetOptionValueById(&option2, variant.OptionValue2)
			if err != nil {
				fmt.Println("Error getting option 2:", err)
				res.ERROR(w, http.StatusInternalServerError, errors.New("Lấy thông tin option 2 thất bại"))
				return
			}
		}

		// Sử dụng thông tin biến thể để tính toán quantity và giá
		effectiveQuantity := cartItem.Quantity
		effectivePrice := float64(effectiveQuantity) * float64(variant.Price)

		cartResponse := response.CartResponse{
			CartID:     cartItem.CartID,
			UserID:     cartItem.UserID,
			ProductID:  variant.ProductID,
			VariantID:  cartItem.VariantID,
			Quantity:   effectiveQuantity,
			TotalPrice: effectivePrice,
			ProductDetail: models.Product{
				ProductID:   product.ProductID,
				Handle:      product.Handle,
				Title:       product.Title,
				Description: product.Description,
				Price:       product.Price,
				CategoryID:  product.CategoryID,
				Image:       product.Image,
				CreatedAt:   product.CreatedAt,
				UpdatedAt:   product.UpdatedAt,
			},
			VariantDetail: variantresponse.VariantDetail{
				VariantID:    variant.VariantID,
				ProductID:    variant.ProductID,
				Title:        variant.Title,
				Price:        variant.Price,
				ComparePrice: variant.ComparePrice,
				CountInStock: variant.CountInStock,
				Image:        variant.Image,
				OptionValue1: models.OptionValue{
					OptionValueID: option1.OptionValueID,
					Value:         option1.Value,
				},
				OptionValue2: models.OptionValue{
					OptionValueID: option2.OptionValueID,
					Value:         option2.Value,
				},
			},
		}

		cartResponses = append(cartResponses, cartResponse)
		totalQuantity += effectiveQuantity
		totalPrice += effectivePrice
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
