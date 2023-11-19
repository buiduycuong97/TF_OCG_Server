package user_handle

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tf_ocg/cmd/app/dbms"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

// delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	uid32 := int32(uid)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var user models.User
	err = dbms.DeleteUser(&user, uid32)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	data := struct {
		Message string `json:"message"`
	}{
		"User deleted successfully",
	}
	res.JSON(w, http.StatusOK, data)
}
