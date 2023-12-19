package discount_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"time"
)

func GetDiscountByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	did, err := strconv.ParseUint(vars["id"], 10, 32)
	discountID := int32(did)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var discount models.Discount
	err = dbms.GetDiscountByID(&discount, discountID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, discount)
}

func GetDiscountByDiscountCode(w http.ResponseWriter, r *http.Request) {
	discountCode := r.URL.Query().Get("discountCode")
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	if discountCode == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Discount code is required"))
		return
	}

	var userDiscount models.UserDiscount
	err = dbms.GetUserDiscountByDiscountCodeAndUserID(&userDiscount, discountCode, int(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var discount models.Discount
			if err := dbms.GetDiscountByDifferentCode(&discount, discountCode); err != nil {
				res.ERROR(w, http.StatusNotFound, errors.New("Invalid discount code"))
				return
			}

			if discount.AvailableQuantity <= 0 {
				res.ERROR(w, http.StatusForbidden, errors.New("Discount has been exhausted"))
				return
			}

			if time.Now().After(discount.EndDate) {
				res.ERROR(w, http.StatusForbidden, errors.New("Discount has expired"))
				return
			}

			res.JSON(w, http.StatusOK, discount)
			return
		}

		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	var discount models.Discount
	if err := dbms.GetDiscountByID(&discount, userDiscount.DiscountID); err != nil {
		res.ERROR(w, http.StatusNotFound, errors.New("Invalid discount code"))
		return
	}
	if discount.AvailableQuantity <= 0 {
		res.ERROR(w, http.StatusForbidden, errors.New("Discount has been exhausted"))
		return
	}

	if time.Now().After(discount.EndDate) {
		res.ERROR(w, http.StatusForbidden, errors.New("Discount has expired"))
		return
	}

	res.JSON(w, http.StatusOK, discount)
}
func GetAllDiscounts(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("searchText")
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

	discounts, err := dbms.GetAllDiscounts(int32(page), int32(pageSize), searchText)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, discounts)
}
