package discount_handle

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"net/http"
	"tf_ocg/cmd/app/dbms"
	"tf_ocg/cmd/app/handler/utils_handle"
	res "tf_ocg/pkg/response_api"
	"tf_ocg/proto/models"
	"time"
)

func CreateDiscountHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Discount models.Discount `json:"discount"`
		UserID   int             `json:"userId"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		res.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	existingDiscount, err := dbms.GetDiscountByCode(requestData.Discount.DiscountCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if err == nil && &existingDiscount != nil {
		res.ERROR(w, http.StatusConflict, errors.New("Discount code already exists"))
		return
	}

	createdDiscount, err := dbms.CreateDiscount(&requestData.Discount)
	if err != nil {
		res.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if requestData.UserID != 0 {
		userDiscount := &models.UserDiscount{
			UserID:     int32(requestData.UserID),
			DiscountID: createdDiscount.DiscountID,
		}
		err := dbms.CreateUserDiscount(userDiscount)
		if err != nil {
			err := dbms.DeleteDiscount(createdDiscount, createdDiscount.DiscountID)
			if err != nil {
				return
			}
			res.ERROR(w, http.StatusInternalServerError, err)
			return
		}
	}

	res.JSON(w, http.StatusCreated, createdDiscount)
}

func CreateAutomaticDiscountForUpgrade(user *models.User) (*models.Discount, error) {
	var discountCode string
	var discountPercentage float64

	switch user.CurrentLevel {
	case models.Silver:
		discountCode = "CHUCMUNGHANGBAC"
		discountPercentage = 5
	case models.Gold:
		discountCode = "CHUCMUNGHANGVANG"
		discountPercentage = 10
	case models.Diamond:
		discountCode = "CHUCMUNGHANGKIMCUONG"
		discountPercentage = 20
	default:
		return nil, errors.New("User is not eligible for an automatic discount at this time")
	}

	automaticDiscount := models.Discount{
		DiscountCode:      discountCode,
		DiscountType:      "percentage",
		Value:             discountPercentage,
		StartDate:         time.Now(),
		EndDate:           time.Now().AddDate(0, 1, 0),
		AvailableQuantity: 1,
	}

	createdDiscount, err := dbms.CreateDiscount(&automaticDiscount)
	if err != nil {
		return nil, err
	}

	userDiscount := models.UserDiscount{
		UserID:     user.UserID,
		DiscountID: createdDiscount.DiscountID,
	}

	err = dbms.CreateUserDiscount(&userDiscount)
	if err != nil {
		return nil, err
	}

	return createdDiscount, nil
}

func GenerateAndSaveDiscountCodes() error {
	discountAmounts := []float64{10000, 20000, 30000, 40000, 50000}
	discountQuantity := 5

	var createdDiscounts []models.Discount

	for _, amount := range discountAmounts {
		discountCode := generateDiscountCode()

		automaticDiscount := models.Discount{
			DiscountCode:      discountCode,
			DiscountType:      "fixed",
			Value:             amount,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 1, 0), // Valid for 1 month
			AvailableQuantity: int32(discountQuantity),
		}

		createdDiscount, err := dbms.CreateDiscount(&automaticDiscount)
		if err != nil {
			return err
		}

		fmt.Printf("Created discount with code %s and amount %f\n", createdDiscount.DiscountCode, createdDiscount.Value)

		createdDiscounts = append(createdDiscounts, *createdDiscount)
	}

	err := utils_handle.SendDiscountMessagesToRabbitMQ(createdDiscounts)
	if err != nil {
		fmt.Printf("Error sending messages to RabbitMQ Cloud: %v\n", err)
	}

	go func(discounts []models.Discount) {
		time.Sleep(2 * time.Minute)
		deleteAllDiscounts(discounts)
	}(createdDiscounts)

	return nil
}

func deleteAllDiscounts(discounts []models.Discount) error {
	for _, discount := range discounts {
		err := dbms.DeleteDiscountAutoGen(&discount, discount.DiscountID)
		if err != nil {
			return err
		}
		fmt.Printf("Deleted discount with code %s and amount %f\n", discount.DiscountCode, discount.Value)
	}
	return nil
}

func generateDiscountCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeLength := 5

	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
