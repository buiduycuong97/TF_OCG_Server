package product_handle

import (
	"encoding/json"
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
		w.Write([]byte("Title is require"))
		return
	}

	var result *models.Product
	result, err = dbms.CreateProduct(&product)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	createProductRes := response.Product{
		ProductID:   result.ProductID,
		Handle:      result.Handle,
		Title:       result.Title,
		Description: result.Description,
		Price:       result.Price,
	}
	res.JSON(w, http.StatusCreated, createProductRes)
}
