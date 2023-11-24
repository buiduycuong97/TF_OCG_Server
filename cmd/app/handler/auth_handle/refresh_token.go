package auth_handle

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing Authorization Header"))
		return
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	claims, _ := utils.VerifyToken(tokenString)
	mapClaims, _ := claims.(jwt.MapClaims)
	idF, ok := mapClaims["id"].(float64)
	if !ok {
		return
	}
	id := int32(idF)
	err := dbms.GetUser(user, id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User not existed"))
		return
	}
	if tokenString != user.RefreshToken {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Fail verify refresh token"))
		return
	}
	accessToken, _ := utils.GenerateAccessToken(id)
	refreshToken, _ := utils.GenerateRefreshToken(id)
	user.RefreshToken = refreshToken
	err = dbms.UpdateUser(user, id)

	tokenRes := RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Update user fail"))
		return
	}
	response_api.JSON(w, http.StatusOK, tokenRes)

}
