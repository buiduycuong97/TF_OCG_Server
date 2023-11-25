package discount_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var discount models.Discount
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &discount)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	createdDiscount, err := dbms.CreateDiscount(&discount)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusCreated, createdDiscount)
}
