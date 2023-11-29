package province_handle

import (
	"errors"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func GetAllProvince(w http.ResponseWriter, r *http.Request) {
	var province []models.Province

	err := dbms.GetAllProvince(&province)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	res.JSON(w, http.StatusOK, province)
}

func GetProvinceByName(w http.ResponseWriter, r *http.Request) {
	provinceName := r.URL.Query().Get("name")

	if provinceName == "" {
		res.ERROR(w, http.StatusBadRequest, errors.New("Province name is required"))
		return
	}

	province, err := dbms.GetProvinceByNameFromDB(provinceName)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, province)
}

func GetProvinceById(w http.ResponseWriter, r *http.Request) {
	provinceIDStr := r.URL.Query().Get("id")

	provinceID, err := strconv.ParseInt(provinceIDStr, 10, 32)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid province ID"))
		return
	}

	province, err := dbms.GetProvinceByIDFromDB(int32(provinceID))
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, province)
}
