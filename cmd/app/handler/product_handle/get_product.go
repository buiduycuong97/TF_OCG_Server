package product_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

// get user by id
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	pid32 := int32(pid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var product models.Product
	err = dbms.GetProductById(&product, pid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, product)
}

func GetListProducts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	products, totalCount, err := dbms.GetListProduct(int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"products":   products,
		"totalPages": calculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, response)
}

func GetListProductByCategoryId(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	categoryIDStr := r.URL.Query().Get("categoryID")

	// Kiểm tra xem categoryID có được cung cấp không
	if categoryIDStr == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("categoryID is required"))
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	products, totalCount, err := dbms.GetListProductByCategoryId(int(categoryID), int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"products":   products,
		"totalPages": calculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, response)
}

func calculateTotalPages(totalCount int64, pageSize int32) int32 {
	totalPages := int32(totalCount) / pageSize
	if int32(totalCount)%pageSize != 0 {
		totalPages++
	}
	return totalPages
}
