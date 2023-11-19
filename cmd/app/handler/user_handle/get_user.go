package user_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	var user []models.User
	err := dbms.GetUsers(&user)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, user)
}

// get user by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	uid32 := int32(uid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}
	var user models.User
	err = dbms.GetUser(&user, uid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res.JSON(w, http.StatusOK, user)
}
