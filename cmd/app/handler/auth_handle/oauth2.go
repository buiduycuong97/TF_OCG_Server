package auth_handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/gomail.v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/dto/user_dto/request"
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
		RedirectURL:  "https://gottago.cyou/api/auth/callback-google",
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
	u := url.URL{
		Host: "gottago.cyou",
		Path: "/login",
	}
	q := make(url.Values)
	q.Set("userID", fmt.Sprintf("%d", result.UserID))
	q.Set("userName", result.UserName)
	q.Set("email", result.Email)
	q.Set("role", result.Role)
	q.Set("phoneNumber", result.PhoneNumber)
	q.Set("userType", result.UserType)
	q.Set("accessToken", accessToken)
	q.Set("refreshToken", refreshToken)
	u.RawQuery = q.Encode()

	urlString := u.String()
	http.Redirect(w, r, urlString, http.StatusTemporaryRedirect)
	fmt.Fprintf(w, "Response : %s", content)
}

func HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	var input request.ForgetPasswordReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user, err := dbms.GetUserByEmail(input.Email)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		response_api.ERROR(w, http.StatusNotFound, errors.New("User not found"))
		return
	}
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user.ResetToken = resetToken
	err = dbms.UpdateUser(user, user.UserID)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	err = sendResetPasswordEmail(user.Email, resetToken)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response_api.JSON(w, http.StatusOK, "Reset password email sent successfully")
}

func HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var input request.ResetPasswordReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response_api.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user, err := dbms.GetUserByResetToken(input.ResetToken)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		response_api.ERROR(w, http.StatusNotFound, errors.New("Invalid reset token"))
		return
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	user.Password = hashedPassword
	user.ResetToken = ""
	err = dbms.UpdateUser(user, user.UserID)
	if err != nil {
		response_api.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response_api.JSON(w, http.StatusOK, "Password reset successfully")
}

func sendResetPasswordEmail(email, resetToken string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	emailAddress := os.Getenv("EMAIL_ADDRESS")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailhost := os.Getenv("EMAIL_HOST")
	subject := "Reset Your Password"
	body := fmt.Sprintf("Click the following link to reset your password: http://localhost:8080/reset-password?resetToken=%s", resetToken)

	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(emailhost, 587, emailAddress, emailPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
