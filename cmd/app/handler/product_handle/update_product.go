package product_handle

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/product_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
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

	productData, err := convertProductToString(product)
	if err != nil {
		log.Println("Lỗi chuyển đổi dữ liệu sản phẩm thành JSON: ", err)
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.SetProductToCache(redisClient, product.Handle, productData)
	if err != nil {
		log.Println("Lưu sản phẩm vào cache thất bại: ", err)
	}

	products, err := dbms.GetListProduct()
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	cacheKey := "list_products"
	err = utils.SetListProductsToCache(redisClient, cacheKey, products)
	if err != nil {
		log.Println("Lưu danh sách sản phẩm vào cache thất bại: ", err)
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
