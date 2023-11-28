package province_handle

import (
	"net/http"
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
