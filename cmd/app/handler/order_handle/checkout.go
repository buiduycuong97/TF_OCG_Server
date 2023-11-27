// Trong package order_handle hoặc một file tương tự
package order_handle

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

	user, err := dbms.GetUserByID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	shippingAddress := r.FormValue("shippingAddress")
	if shippingAddress == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Shipping address is required"))
		return
	}

	provinceIDStr := r.FormValue("provinceId")
	provinceID, err := strconv.Atoi(provinceIDStr)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid province ID format"))
		return
	}

	cartItems, err := dbms.GetCartByUserID(userID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(cartItems) == 0 {
		res.ERROR(w, http.StatusBadRequest, errors.New("Cart is empty, cannot proceed with checkout"))
		return
	}

	discountCode := r.FormValue("discountCode")
	var discountAmount float64
	if discountCode != "" {
		discount, err := dbms.GetDiscountByCode(database.Db, discountCode)

		usedDiscount, err := dbms.CheckUserUsedDiscount(userID, discount.DiscountID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to check if user used discount: %v", err))
			return
		}

		if !usedDiscount {
			res.ERROR(w, http.StatusBadRequest, errors.New("Discount code has already been used by this user"))
			return
		}

		discountAmount, err = dbms.ApplyDiscountForOrder(database.Db, cartItems, discountCode)
		if err != nil {
			res.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = dbms.MarkDiscountAsUsed(userID, discount.DiscountID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to mark discount as used: %v", err))
			return
		}
	}

	newOrder := models.Order{
		UserID:          userID,
		OrderDate:       time.Now(),
		ShippingAddress: shippingAddress,
		Status:          models.Pending,
		ProvinceID:      int32(provinceID),
		TotalQuantity:   0,
		TotalPrice:      0.0,
		DiscountAmount:  discountAmount,
	}

	tx := database.Db.Begin()

	createdOrder, err := dbms.CreateOrder(tx, &newOrder)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Created Order ID:", createdOrder.OrderID)

	province, err := dbms.GetProvinceByID(newOrder.ProvinceID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

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

		newOrder.TotalQuantity += cartItem.Quantity
		newOrder.TotalPrice += cartItem.TotalPrice
	}

	newOrder.TotalPrice -= newOrder.DiscountAmount

	err = dbms.ClearCart(tx, userID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to clear cart: %v", err))
		return
	}

	err = dbms.UpdateOrderTotalValues(tx, createdOrder.OrderID, newOrder.TotalQuantity, newOrder.TotalPrice)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to update order values: %v", err))
		return
	}

	tx.Commit()

	responseData := map[string]interface{}{
		"message":          "Order created successfully",
		"user_name":        user.UserName,
		"phone_number":     user.PhoneNumber,
		"province_name":    province.ProvinceName,
		"shipping_fee":     province.ShippingFee,
		"order_id":         createdOrder.OrderID,
		"order_date":       createdOrder.OrderDate,
		"shipping_address": createdOrder.ShippingAddress,
		"total_quantity":   newOrder.TotalQuantity,
		"total_price":      newOrder.TotalPrice,
		"discount_amount":  newOrder.DiscountAmount,
	}

	res.JSON(w, http.StatusCreated, responseData)
}
