package order_handle

import (
	"errors"
	"fmt"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	database "tf_ocg/pkg/database_manager"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"time"
)

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	shippingAddress := r.FormValue("shipping_address")
	if shippingAddress == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Shipping addres5s is required"))
		return
	}

	cartItems, err := dbms.GetCartByUserID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	newOrder := models.Order{
		UserID:          userID,
		OrderDate:       time.Now(),
		ShippingAddress: shippingAddress,
		Status:          models.Pending,
	}

	tx := database.Db.Begin()

	createdOrder, err := dbms.CreateOrder(tx, &newOrder)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Created Order ID:", createdOrder.OrderID)

	for _, cartItem := range cartItems {
		orderDetail := &models.OrderDetail{
			OrderID:   createdOrder.OrderID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.TotalPrice / float64(cartItem.Quantity),
		}

		err = dbms.CreateOrderDetail(tx, orderDetail)
		if err != nil {
			tx.Rollback()
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	err = dbms.ClearCart(tx, userID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to clear cart: %v", err))
		return
	}

	tx.Commit()

	res.JSON(w, http.StatusCreated, map[string]string{"message": "Order created successfully"})
}
