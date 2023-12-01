package product_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/product_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	pid32 := int32(pid)

	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var product models.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = dbms.UpdateProduct(&product, pid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	productResponse := response.ProductResponseUpdate{
		ProductID:   product.ProductID,
		Handle:      product.Handle,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		UpdatedAt:   product.UpdatedAt,
	}

	res.JSON(w, http.StatusOK, productResponse)
}

//
//func UpdateProductQuantityHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	productID, err := strconv.ParseInt(vars["id"], 10, 32)
//	if err != nil {
//		res.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//
//	quantity, err := strconv.Atoi(r.FormValue("quantity"))
//	if err != nil {
//		res.ERROR(w, http.StatusBadRequest, err)
//		return
//	}
//
//	err = dbms.UpdateProductQuantity(int32(productID), int32(quantity))
//	if err != nil {
//		res.ERROR(w, http.StatusInternalServerError, err)
//		return
//	}
//
//	res.JSON(w, http.StatusOK, map[string]string{"message": "Product quantity updated successfully"})
//}
