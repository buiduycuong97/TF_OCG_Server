package auth_handle

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	dto "tf_ocg/cmd/app/dto/user_dto"
	"tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"tf_ocg/utils"
)

type googleResponse struct {
	Name  string
	Email string
}

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/callback-google",
		ClientID:     "275944160808-8b24hbsrsodun2vtd1ubobih7ll1bflm.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
	randomStateGoogle = "random"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/auth/login-google">Google</a></body></html`
	fmt.Fprint(w, html)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomStateGoogle)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomStateGoogle {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	NoContext := context.TODO()
	token, err := googleOauthConfig.Exchange(NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token &s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create token &s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not parse response &s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var googleRes googleResponse
	err = json.Unmarshal(content, &googleRes)
	if err != nil {
		response_api.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//create user
	var user models.User
	user.UserName = googleRes.Name
	user.Email = googleRes.Email
	user.UserType = "google"

	var result *models.User
	result, err = dbms.CreateUser(&user)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	// save rfToken to db
	accessToken, _ := utils.GenerateAccessToken(user.UserID)
	refreshToken, _ := utils.GenerateRefreshToken(user.UserID)
	user.RefreshToken = refreshToken
	err = dbms.UpdateUser(&user, user.UserID)
	loginRes := dto.LoginRes{
		UserID:       result.UserID,
		UserName:     result.UserName,
		Email:        result.Email,
		Role:         result.Role,
		UserType:     result.UserType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	response_api.JSON(w, http.StatusOK, loginRes)
	fmt.Fprintf(w, "Response : %s", content)
}