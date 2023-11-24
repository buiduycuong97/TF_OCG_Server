package review_handle

import (
	"encoding/json"
	"errors"
	goaway "github.com/TwiN/go-away"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
)

func AddReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		res.ERROR(w, http.StatusUnsupportedMediaType, errors.New("Content-Type must be application/json"))
		return
	}

	// Giải mã JSON từ body của request vào struct Review
	var newReview models.Review
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newReview)
	if err != nil {
		res.ERROR(w, http.StatusBadRequest, errors.New("Invalid JSON format in request body"))
		return
	}

	// Lấy userID từ request header
	userID, err := utils_handle.GetUserIDFromRequest(r)
	if err != nil {
		res.ERROR(w, http.StatusUnauthorized, errors.New("Invalid token"))
		return
	}

	// Gán giá trị UserID từ request header vào struct Review
	newReview.UserID = userID

	// Nếu chỉ có productID, chỉ trả về danh sách
	if newReview.Rating == 0 && newReview.Comment == "" {
		reviews, err := dbms.GetReviewsByProductID(newReview.ProductID)
		if err != nil {
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		res.JSON(w, http.StatusOK, reviews)
		return
	}

	// Kiểm tra sensitive trước khi thêm review
	if isSensitive(newReview.Comment) {
		sendSensitiveNotification(userID, "Your review has been flagged as sensitive. It will not be published.")
		res.ERROR(w, http.StatusForbidden, errors.New("Sensitive content detected"))
		return
	}

	// Thực hiện thêm review vào cơ sở dữ liệu
	err = dbms.CreateReview(&newReview)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Lấy danh sách reviews theo ProductID từ cơ sở dữ liệu
	reviews, err := dbms.GetReviewsByProductID(newReview.ProductID)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	// Trả về danh sách reviews dưới dạng JSON
	res.JSON(w, http.StatusOK, reviews)
}

// isSensitive kiểm tra xem có nội dung nhạy cảm không
func isSensitive(comment string) bool {
	isProfane := goaway.IsProfane(comment)
	log.Println(isProfane)
	return isProfane
}

func sendSensitiveNotification(userID int32, message string) error {
	email, err := getEmailByUserID(userID)
	if err != nil {
		return err
	}
	emailAddress := "pau30012002@gmail.com"
	emailPassword := "pljf fqgx yycq ynhq"
	subject := "Sensitive Content Notification"
	body := message

	m := gomail.NewMessage()
	m.SetHeader("From", emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, emailAddress, emailPassword)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// getEmailByUserID là một hàm giả sử để lấy địa chỉ email từ UserID
func getEmailByUserID(userID int32) (string, error) {
	email, err := dbms.GetEmailByUserID(userID)
	if err != nil {
		return "", err
	}
	return email, nil
}
