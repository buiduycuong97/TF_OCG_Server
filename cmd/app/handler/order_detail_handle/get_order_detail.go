package order_detail_handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/order_detail/response"
	responseproduct "tf_ocg/cmd/app/dto/product_dto/response"
	responseuser "tf_ocg/cmd/app/dto/user_dto/response"
	responsevariant "tf_ocg/cmd/app/dto/variant_dto/response"
)

func GetOrderInfo(orderID int32) (response.OrderInfo, error) {
	order, err := dbms.GetOrderByID(orderID)
	if err != nil {
		return response.OrderInfo{}, err
	}

	orderDetails, err := dbms.GetOrderDetailsByOrderID(orderID)
	if err != nil {
		return response.OrderInfo{}, err
	}

	var orderDetailsInfo []response.OrderDetailInfo
	for _, orderDetail := range orderDetails {
		user, err := dbms.GetUserByID(order.UserID)
		if err != nil {
			return response.OrderInfo{}, err
		}

		variant, err := dbms.GetVariantByIdInGetOrder(orderDetail.VariantID)
		if err != nil {
			return response.OrderInfo{}, err
		}

		product, err := dbms.GetProductByID(variant.ProductID)
		if err != nil {
			return response.OrderInfo{}, err
		}

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
				Image:        variant.Image,
				OptionValue1: variant.OptionValue1,
				OptionValue2: variant.OptionValue2,
			},
		}

		orderDetailsInfo = append(orderDetailsInfo, orderDetailInfo)
	}

	orderInfo := response.OrderInfo{
		OrderID:         order.OrderID,
		UserID:          order.UserID,
		OrderDate:       order.OrderDate,
		ShippingAddress: order.ShippingAddress,
		PhoneOrder:      order.PhoneOrder,
		Status:          order.Status,
		ProvinceID:      order.ProvinceID,
		TotalQuantity:   order.TotalQuantity,
		TotalPrice:      order.TotalPrice,
		DiscountAmount:  order.DiscountAmount,
		GrandTotal:      order.GrandTotal,
		OrderDetails:    orderDetailsInfo,
		UpdatedAt:       order.UpdatedAt,
	}

	return orderInfo, nil
}

func GetOrderInfoHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("orderID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid orderID", http.StatusBadRequest)
		return
	}

	orderInfo, err := GetOrderInfo(int32(orderID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving order information: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderInfo)
}
