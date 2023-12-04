package order_handle

import (
	"errors"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func RequestCancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("orderId")
	if orderIDStr == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("orderId is required"))
		return
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid orderId format"))
		return
	}

	currentStatus, err := dbms.GetOrderStatus(int32(orderID))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if currentStatus != string(models.Pending) {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid status transition"))
		return
	}

	err = dbms.UpdateOrderStatus(int32(orderID), string(models.Cancelled))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Order cancellation requested successfully"})
}
