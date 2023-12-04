package auth_handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/user_dto/request"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//read from body
	var body request.LoginReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// validate
	if (body.Email) == "" || (body.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password to obtain the token"))
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: body.Password,
	}
	var userRes *models.User
	fmt.Println(user)
	userRes, err = dbms.LoginUser(&user)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := utils.GenerateAccessToken(userRes.UserID)
	refreshToken, err := utils.GenerateRefreshToken(userRes.UserID)

	// save rfToken to db
	userRes.RefreshToken = refreshToken
	err = dbms.UpdateUser(userRes, userRes.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating JWT token"))
	}

	// set cookie
	cookie1 := &http.Cookie{
		Name:    "accessToken",
		Path:    "/",
		Value:   accessToken,
		Expires: time.Now().Add(time.Hour * 1),
	}
	http.SetCookie(w, cookie1)
	cookie2 := &http.Cookie{
		Name:    "refreshToken",
		Path:    "/",
		Value:   refreshToken,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, cookie2)

	loginRes := request.LoginRes{
		UserID:       userRes.UserID,
		UserName:     userRes.UserName,
		Email:        userRes.Email,
		PhoneNumber:  userRes.PhoneNumber,
		Role:         userRes.Role,
		UserType:     userRes.UserType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Authorization", "Bearer "+accessToken)
	w.Header().Set("user", fmt.Sprintf("%+v", loginRes)) //test
	res.JSON(w, http.StatusOK, loginRes)

}

// login admin
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	//read from body
	var body request.LoginReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// validate
	if (body.Email) == "" || (body.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password to obtain the token"))
		return
	}

	user := models.User{
		Email:    body.Email,
		Password: body.Password,
	}
	var userRes *models.User
	fmt.Printf("%+v", user)
	userRes, err = dbms.LoginAdmin(&user)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := utils.GenerateAccessTokenAdmin(userRes.Role)
	refreshToken, err := utils.GenerateRefreshTokenAdmin(userRes.Role)

	// save rfToken to db
	userRes.RefreshToken = refreshToken
	err = dbms.UpdateUser(userRes, userRes.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating JWT token"))
	}

	// set cookie
	cookie := &http.Cookie{
		Name:    "admin",
		Path:    "/",
		Value:   accessToken,
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, cookie)

	loginRes := request.LoginRes{
		UserID:       userRes.UserID,
		UserName:     userRes.UserName,
		Email:        userRes.Email,
		Role:         userRes.Role,
		UserType:     userRes.UserType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Authorization", "Bearer "+accessToken)
	w.Header().Set("user", fmt.Sprintf("%+v", loginRes))
	res.JSON(w, http.StatusOK, loginRes)

}
