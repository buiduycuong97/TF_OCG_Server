package user_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/user_dto"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//validate
	if user.UserName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name is require"))
		return
	}
	if !(utils.IsValidEmail(user.Email)) || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is not valid"))
		return
	}

	var result *models.User
	result, err = dbms.CreateUser(&user)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	createUserRes := user_dto.CreateUserReq{
		UserID:   result.UserID,
		UserName: result.UserName,
		Email:    result.Email,
		Role:     result.Role,
		UserType: result.UserType,
	}
	res.JSON(w, http.StatusCreated, createUserRes)
}