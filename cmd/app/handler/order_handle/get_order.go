package order_handle

import (
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	order_response "tf_ocg/cmd/app/dto/order/response"
	order_detail_response "tf_ocg/cmd/app/dto/order_detail/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func ViewOrderHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := dbms.GetOrdersByStatus(models.CompleteTheOrder)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to get orders"))
		return
	}

	var orderResponses []order_response.OrderResponse
	for _, order := range orders {
		orderDetails, err := dbms.GetOrderDetailsByOrderID(order.OrderID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to get order details"))
			return
		}

		orderResponse := order_response.OrderResponse{
			OrderID:         order.OrderID,
			UserID:          order.UserID,
			OrderDate:       order.OrderDate,
			ShippingAddress: order.ShippingAddress,
			Status:          order.Status,
			OrderDetails:    make([]order_detail_response.OrderDetailResponse, 0),
		}

		for _, orderDetail := range orderDetails {
			orderDetailResponse := order_detail_response.OrderDetailResponse{
				OrderDetailID: orderDetail.OrderDetailID,
				ProductID:     orderDetail.ProductID,
				Quantity:      orderDetail.Quantity,
				Price:         orderDetail.Price,
			}

			orderResponse.OrderDetails = append(orderResponse.OrderDetails, orderDetailResponse)
		}

		orderResponses = append(orderResponses, orderResponse)
	}

	res.JSON(w, http.StatusOK, orderResponses)
}
