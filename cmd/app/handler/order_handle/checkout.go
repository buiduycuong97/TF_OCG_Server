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

	// Lấy thông tin người dùng từ cơ sở dữ liệu
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

	newOrder := models.Order{
		UserID:          userID,
		OrderDate:       time.Now(),
		ShippingAddress: shippingAddress,
		Status:          models.Pending,
		ProvinceID:      int32(provinceID),
	}

	tx := database.Db.Begin()

	createdOrder, err := dbms.CreateOrder(tx, &newOrder)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println("Created Order ID:", createdOrder.OrderID)

	// Lấy thông tin về tỉnh từ cơ sở dữ liệu
	province, err := dbms.GetProvinceByID(newOrder.ProvinceID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Tạo một map chứa thông tin cần trả về
	responseData := map[string]interface{}{
		"message":          "Order created successfully",
		"user_name":        user.UserName,
		"phone_number":     user.PhoneNumber,
		"province_name":    province.ProvinceName,
		"shipping_fee":     province.ShippingFee,
		"order_id":         createdOrder.OrderID,
		"order_date":       createdOrder.OrderDate,
		"shipping_address": createdOrder.ShippingAddress,
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
	}

	err = dbms.ClearCart(tx, userID)
	if err != nil {
		tx.Rollback()
		res.ERROR(w, http.StatusInternalServerError, fmt.Errorf("Failed to clear cart: %v", err))
		return
	}

	tx.Commit()

	res.JSON(w, http.StatusCreated, responseData)
}
