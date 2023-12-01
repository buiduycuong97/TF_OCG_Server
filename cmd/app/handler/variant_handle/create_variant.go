package variant_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/variant_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CreateVariantHandler(w http.ResponseWriter, r *http.Request) {
	var variant models.Variant

	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &variant)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Kiểm tra xem variant có hợp lệ không
	if variant.ProductID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ProductID is required"))
		return
	}

	if variant.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Title is required"))
		return
	}

	// Tạo variant trong cơ sở dữ liệu
	createdVariant, err := dbms.CreateVariant(&variant)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Trả về thông tin của variant vừa được tạo
	createVariantRes := response.VariantResponse{
		VariantID:      createdVariant.VariantID,
		ProductID:      createdVariant.ProductID,
		Title:          createdVariant.Title,
		Price:          createdVariant.Price,
		ComparePrice:   createdVariant.ComparePrice,
		CountInStock:   createdVariant.CountInStock,
		OptionProduct1: createdVariant.OptionProduct1,
		OptionProduct2: createdVariant.OptionProduct2,
	}

	res.JSON(w, http.StatusCreated, createVariantRes)
}
