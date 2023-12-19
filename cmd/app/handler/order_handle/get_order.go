package order_handle

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	order_response "tf_ocg/cmd/app/dto/order/response"
	order_detail_response "tf_ocg/cmd/app/dto/order_detail/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func ViewPendingOrdersHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.Pending)
}
func ViewOrderBeingDeliveredHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.OrderBeingDelivered)
}
func ViewCompleteTheOrderHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.CompleteTheOrder)
}
func ViewCancelledOrdersHandler(w http.ResponseWriter, r *http.Request) {
	viewOrdersByStatus(w, r, models.Cancelled)
}

func viewOrdersByStatus(w http.ResponseWriter, r *http.Request, status models.OrderStatus) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid or missing 'page' parameter"))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil || pageSize < 1 {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid or missing 'pageSize' parameter"))
		return
	}

	orders, totalItem, err := dbms.GetOrdersByStatus(status, page, pageSize)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	totalPages := int64(0)
	if pageSize > 0 {
		totalPages = (totalItem + pageSize - 1) / pageSize
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
				VariantImage:  getImageByVariantID(orderDetail.VariantID),
			}

			orderResponse.OrderDetails = append(orderResponse.OrderDetails, orderDetailResponse)
		}

		orderResponses = append(orderResponses, orderResponse)
	}

	response := map[string]interface{}{
		"orders":     orderResponses,
		"totalItem":  totalItem,
		"totalPages": totalPages,
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

func GetAllOrder(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	if status == "pending" {
		status = "pending"
	} else if status == "onBeingDelivered" {
		status = "order being delivered"
	} else if status == "completeTheOrder" {
		status = "complete the order"
	} else if status == "cancelled" {
		status = "cancelled"
	} else {
		status = ""
	}
	orders, err := dbms.GetAllOrder(int32(page), int32(pageSize), status)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, orders)

}

func getImageByVariantID(variantID int32) string {
	image, err := dbms.GetImageByVariantID(variantID)
	if err != nil {
		fmt.Printf("Error getting image for variant ID %d: %v\n", variantID, err)
		return ""
	}
	return image
}
