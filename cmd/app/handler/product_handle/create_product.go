package product_handle

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/product_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &product)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if product.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Title is required"))
		return
	}
	product.Handle = slug.Make(product.Title)

	_, err = dbms.CreateProduct(&product)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	createProductRes := response.Product{
		ProductId:         product.ProductID,
		Handle:            product.Handle,
		Title:             product.Title,
		Description:       product.Description,
		Price:             product.Price,
		CategoryID:        product.CategoryID,
		QuantityRemaining: product.QuantityRemaining,
	}
	res.JSON(w, http.StatusCreated, createProductRes)
}
