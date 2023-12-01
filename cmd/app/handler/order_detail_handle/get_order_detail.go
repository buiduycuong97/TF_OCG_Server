package order_detail_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/order_detail/response"
	responseproduct "tf_ocg/cmd/app/dto/product_dto/response"
	responseuser "tf_ocg/cmd/app/dto/user_dto/response"
	responsevariant "tf_ocg/cmd/app/dto/variant_dto/response"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
)

func GetOrderInfo(orderID int32, userID int32) (response.OrderInfo, error) {
	// Lấy thông tin order
	order, err := dbms.GetOrderByID(orderID)
	if err != nil {
		return response.OrderInfo{}, err
	}

	// Lấy danh sách order details
	orderDetails, err := dbms.GetOrderDetailsByOrderID(orderID)
	if err != nil {
		return response.OrderInfo{}, err
	}

	// Tạo slice để chứa thông tin của từng order detail
	var orderDetailsInfo []response.OrderDetailInfo
	// Duyệt qua từng order detail để lấy thông tin user và product
	for _, orderDetail := range orderDetails {
		// Lấy thông tin user từ orderDetail.UserID
		user, err := dbms.GetUserByID(userID)
		if err != nil {
			return response.OrderInfo{}, err
		}

		variant, err := dbms.GetVariantByIdInGetOrder(orderDetail.VariantID)
		if err != nil {
			return response.OrderInfo{}, err
		}

		// Lấy thông tin product từ orderDetail.ProductID
		product, err := dbms.GetProductByID(variant.ProductID)
		if err != nil {
			return response.OrderInfo{}, err
		}

		// Tạo struct mới chứa thông tin
		orderDetailInfo := response.OrderDetailInfo{
			OrderDetailID: orderDetail.OrderDetailID,
			OrderID:       orderDetail.OrderID,
			Quantity:      orderDetail.Quantity,
			Price:         orderDetail.Price,
			UserInfo: responseuser.UserInfo{
				UserID:      user.UserID,
				UserName:    user.UserName,
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
			},
			ProductInfo: responseproduct.ProductInfo{
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
			VariantInfo: responsevariant.VariantResponse{
				VariantID:    variant.VariantID,
				Title:        variant.Title,
				Price:        variant.Price,
				ComparePrice: variant.ComparePrice,
				CountInStock: variant.CountInStock,
				OptionValue1: variant.OptionValue1,
				OptionValue2: variant.OptionValue2,
			},
		}

		orderDetailsInfo = append(orderDetailsInfo, orderDetailInfo)
	}

	// Tạo struct chứa thông tin order và danh sách order details
	orderInfo := response.OrderInfo{
		OrderID:         order.OrderID,
		UserID:          order.UserID,
		OrderDate:       order.OrderDate,
		ShippingAddress: order.ShippingAddress,
		Status:          order.Status,
		ProvinceID:      order.ProvinceID,
		TotalQuantity:   order.TotalQuantity,
		TotalPrice:      order.TotalPrice,
		DiscountAmount:  order.DiscountAmount,
		GrandTotal:      order.GrandTotal,
		OrderDetails:    orderDetailsInfo,
	}

	return orderInfo, nil
}

func GetOrderInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy orderID từ URL parameter
	orderIDStr := r.URL.Query().Get("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid orderID", http.StatusBadRequest)
		return
	}

	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	// Gọi hàm GetOrderInfo để lấy thông tin chi tiết về order và orderDetails
	orderInfo, err := GetOrderInfo(int32(orderID), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving order information: %v", err), http.StatusInternalServerError)
		return
	}

	// Chuyển đổi struct orderInfo thành JSON và gửi về client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderInfo)
}
