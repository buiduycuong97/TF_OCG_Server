package user_handle

import (
	"encoding/json"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/user_dto/response"
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
	if user.PhoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Phone is require"))
		return
	}
	if !(utils.IsValidEmail(user.Email)) || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is not valid"))
		return
	}
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password is require"))
		return
	}

	var result *models.User
	result, err = dbms.CreateUser(&user)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	createUserRes := response.CreateUserRes{
		UserID:      result.UserID,
		UserName:    result.UserName,
		Email:       result.Email,
		PhoneNumber: result.PhoneNumber,
		Role:        result.Role,
		UserType:    result.UserType,
	}
	res.JSON(w, http.StatusCreated, createUserRes)
}
