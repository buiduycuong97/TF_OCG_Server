package dbms

import (
	"errors"
	database "tf_ocg/pkg/database_manager"
	"tf_ocg/utils"
)
import "tf_ocg/proto/models"

// create a user
func CreateUser(user *models.User) (*models.User, error) {
	newUser := models.User{}

	database.Db.Raw("SELECT * FROM users WHERE email = ? ", user.Email).Scan(&newUser)
	if user.Email == newUser.Email {
		return nil, errors.New("Email existed!")
	} else {
		hashPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashPassword

		err = database.Db.Create(user).Error
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

// get users
func GetUsers(user *[]models.User) (err error) {
	err = database.Db.Find(user).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by id
func GetUser(user *models.User, id int32) (err error) {
	err = database.Db.Where("user_id = ?", id).First(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByID(userID int32) (*models.User, error) {
	var user models.User
	err := database.Db.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// update user
func UpdateUser(user *models.User, id int32) (err error) {
	database.Db.Model(user).Where("user_id = ?", id).Updates(user)
	return nil
}

// delete user
func DeleteUser(user *models.User, userID int32) error {
	if err := DeleteUserCart(userID); err != nil {
		return err
	}
	if err := DeleteUserReviews(userID); err != nil {
		return err
	}
	if err := DeleteUserOrders(userID); err != nil {
		return err
	}
	if err := DeleteUserDiscounts(userID); err != nil {
		return err
	}
	return database.Db.Delete(user, "user_id = ?", userID).Error
}

func DeleteUserCart(userID int32) error {
	return database.Db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error
}

func DeleteUserReviews(userID int32) error {
	return database.Db.Where("user_id = ?", userID).Delete(&models.Review{}).Error
}

func DeleteUserOrders(userID int32) error {
	var orders []models.Order
	if err := database.Db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return err
	}

	for _, order := range orders {
		if err := DeleteOrder(order.OrderID); err != nil {
			return err
		}
	}

	return nil
}

func DeleteOrder(orderID int32) error {
	if err := database.Db.Where("order_id = ?", orderID).Delete(&models.OrderDetail{}).Error; err != nil {
		return err
	}
	return database.Db.Where("order_id = ?", orderID).Delete(&models.Order{}).Error
}

func DeleteUserDiscounts(userID int32) error {
	return database.Db.Where("user_id = ?", userID).Delete(&models.UserDiscount{}).Error
}

// login user
func LoginUser(user *models.User) (*models.User, error) {
	userRes := &models.User{}

	database.Db.Raw("SELECT * FROM users WHERE email = ?", user.Email).Scan(userRes)
	if userRes.Role != "user" {
		return nil, errors.New("You are not user")
	}
	match := utils.CheckPasswordHash(user.Password, userRes.Password)
	if !match {
		return nil, errors.New("Wrong password")
	} else {
		return userRes, nil
	}
}

// login admin
func LoginAdmin(user *models.User) (userRes *models.User, err error) {
	userRes = &models.User{}
	database.Db.Raw("SELECT * FROM users WHERE email = ?", user.Email).Scan(userRes)
	if userRes.Role != "admin" {
		return nil, errors.New("You dont have permission")
	}
	match := utils.CheckPasswordHash(user.Password, userRes.Password)
	if !match {
		return nil, errors.New("Wrong password")
	} else {
		return userRes, nil
	}
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := database.Db.Where("email = ?", email).First(user).Error
	return user, err
}

func GetUserByResetToken(resetToken string) (*models.User, error) {
	var user models.User
	err := database.Db.Where("reset_token = ?", resetToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetEmailByUserID(userID int32) (string, error) {
	var user models.User
	err := database.Db.Select("email").Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

func SearchUser(searchText string, page, pageSize int32) ([]*models.User, error) {
	offset := (page - 1) * pageSize
	users := []*models.User{}

	query := database.Db

	if searchText != "" {
		query = query.Where("user_name LIKE ?", "%"+searchText+"%")
	}

	query = query.Offset(int(offset)).Limit(int(pageSize))

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func ChangePassword(email string, password string) error {
	var user models.User
	hashPass, _ := utils.HashPassword(password)
	database.Db.Model(&user).Where("email = ?", email).Update("password", hashPass)
	return nil
}
