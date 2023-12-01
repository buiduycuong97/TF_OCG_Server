package order_handle

import (
	"errors"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	order_response "tf_ocg/cmd/app/dto/order/response"
	order_detail_response "tf_ocg/cmd/app/dto/order_detail/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

//	func ViewOrderHandler(w http.ResponseWriter, r *http.Request) {
//		orders, err := dbms.GetOrdersByStatus(models.CompleteTheOrder)
//		if err != nil {
//			res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to get orders"))
//			return
//		}
//
//		var orderResponses []order_response.OrderResponse
//		for _, order := range orders {
//			orderDetails, err := dbms.GetOrderDetailsByOrderID(order.OrderID)
//			if err != nil {
//				res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to get order details"))
//				return
//			}
//
//			orderResponse := order_response.OrderResponse{
//				OrderID:         order.OrderID,
//				UserID:          order.UserID,
//				OrderDate:       order.OrderDate,
//				ShippingAddress: order.ShippingAddress,
//				Status:          order.Status,
//				OrderDetails:    make([]order_detail_response.OrderDetailResponse, 0),
//			}
//
//			for _, orderDetail := range orderDetails {
//				orderDetailResponse := order_detail_response.OrderDetailResponse{
//					OrderDetailID: orderDetail.OrderDetailID,
//					ProductID:     orderDetail.ProductID,
//					Quantity:      orderDetail.Quantity,
//					Price:         orderDetail.Price,
//				}
//
//				orderResponse.OrderDetails = append(orderResponse.OrderDetails, orderDetailResponse)
//			}
//
//			orderResponses = append(orderResponses, orderResponse)
//		}
//
//		res.JSON(w, http.StatusOK, orderResponses)
//	}
func ViewPendingOrdersHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.Pending)
}
func ViewOrderBeingDeliveredHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.OrderBeingDelivered)
}
func ViewCompleteTheOrderHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.CompleteTheOrder)
}
func ViewRequestToCancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.RequestToCancelOrder)
}
func ViewCancelledOrdersHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.Cancelled)
}

func viewOrdersByStatus(w http.ResponseWriter, r *http.Request, status models.OrderStatus) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	defaultPage := int32(1)
	defaultPageSize := int32(4)

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = int64(defaultPage)
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = int64(defaultPageSize)
	}

	orders, totalItem, err := dbms.GetOrdersByStatus(status, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var orderResponses []order_response.OrderResponseList
	for _, order := range orders {
		orderDetails, err := dbms.GetOrderDetailsByOrderID(order.OrderID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, errors.New("Failed to get order details"))
			return
		}

		orderResponse := order_response.OrderResponseList{
			OrderID:         order.OrderID,
			UserID:          order.UserID,
			OrderDate:       order.OrderDate,
			ShippingAddress: order.ShippingAddress,
			Status:          order.Status,
			OrderDetails:    make([]order_detail_response.OrderDetailResponse, 0),
			TotalPrice:      calculateTotalPrice(orderDetails),
		}

		for _, orderDetail := range orderDetails {
			orderDetailResponse := order_detail_response.OrderDetailResponse{
				OrderDetailID: orderDetail.OrderDetailID,
				VariantID:     orderDetail.VariantID,
				Quantity:      orderDetail.Quantity,
				Price:         orderDetail.Price,
			}

			orderResponse.OrderDetails = append(orderResponse.OrderDetails, orderDetailResponse)
		}

		orderResponses = append(orderResponses, orderResponse)
	}

	response := map[string]interface{}{
		"orders":    orderResponses,
		"totalItem": totalItem,
	}

	res.JSON(w, http.StatusOK, response)
}

func calculateTotalPrice(orderDetails []models.OrderDetail) float64 {
	var totalPrice float64

	for _, orderDetail := range orderDetails {
		totalPrice += float64(orderDetail.Quantity) * orderDetail.Price
	}

	return totalPrice
}
