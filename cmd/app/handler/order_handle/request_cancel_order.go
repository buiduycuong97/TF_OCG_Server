package order_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func RequestCancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Đọc dữ liệu JSON từ yêu cầu
	var requestData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Error decoding JSON"))
		return
	}

	// Trích xuất orderId từ dữ liệu JSON
	orderIDRaw, ok := requestData["orderId"]
	if !ok {
		res.ERROR(w, http.StatusBadRequest, errors.New("orderId is required"))
		return
	}

	// Kiểm tra định dạng của orderId
	orderID, err := strconv.ParseInt(fmt.Sprintf("%v", orderIDRaw), 10, 32)
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
