package order_handle

import (
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func AdminAcceptCancelRequestHandler(w http.ResponseWriter, r *http.Request) {
	orderID, err := getOrderIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	currentStatus, err := dbms.GetOrderStatus(orderID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if currentStatus != string(models.RequestToCancelOrder) {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid status transition"))
		return
	}

	err = dbms.UpdateOrderStatus(orderID, string(models.Cancelled))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	orderDetails, err := dbms.GetOrderDetailsByOrderID(orderID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	for _, orderItem := range orderDetails {
		err := dbms.UpdateProductQuantityWithIncrease(orderItem.ProductID, orderItem.Quantity)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Order cancellation accepted by admin"})
}
