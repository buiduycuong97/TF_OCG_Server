package category_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/category_dto/response"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Categories
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &category)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if category.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is require"))
		return
	}

	var result *models.Categories
	result, err = dbms.CreateCategory(&category)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	createCategoryRes := response.CategoryResponse{
		Name:   result.Name,
		Handle: result.Handle,
		Image:  result.Image,
	}
	res.JSON(w, http.StatusCreated, createCategoryRes)
}
