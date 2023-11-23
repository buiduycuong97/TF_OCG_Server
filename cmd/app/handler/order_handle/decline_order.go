package order_handle

import (
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func AdminDeclineOrderHandler(w http.ResponseWriter, r *http.Request) {
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

	if currentStatus != string(models.Pending) {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid status transition"))
		return
	}

	err = dbms.UpdateOrderStatus(orderID, string(models.Cancelled))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Order decline by admin"})
}
