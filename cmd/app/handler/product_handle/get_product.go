package product_handle

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
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

func GetProductByHandle(w http.ResponseWriter, r *http.Request) {
	handle := r.URL.Query().Get("handle")
	var product models.Product
	err := dbms.GetProductByHandle(&product, handle)
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
		"totalPages": utils_handle.CalculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, response)
}

func GetListProductByCategoryId(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	categoryIDStr := r.URL.Query().Get("categoryId")

	categoryId, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		http.Error(w, "Invalid categoryId", http.StatusBadRequest)
		return
	}
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

	products, totalCount, err := dbms.GetListProductByCategoryId(categoryId, int32(page), int32(pageSize))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"products":   products,
		"totalPages": utils_handle.CalculateTotalPages(totalCount, int32(pageSize)),
	}

	res.JSON(w, http.StatusOK, response)
}
