package order_handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CompleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		OrderID int32 `json:"orderId"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Failed to decode request body"))
		return
	}

	orderID := requestBody.OrderID

	currentStatus, err := dbms.GetOrderStatus(orderID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if currentStatus != string(models.OrderBeingDelivered) {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid status transition"))
		return
	}

	err = dbms.UpdateOrderStatus(orderID, string(models.CompleteTheOrder))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"message": "Order completed successfully"})
}
