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
func GetUsers(User *[]models.User) (err error) {
	err = database.Db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by id
func GetUser(User *models.User, id int32) (err error) {
	err = database.Db.Where("user_id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

// update user
func UpdateUser(User *models.User, id int32) (err error) {
	database.Db.Model(User).Where("user_id = ?", id).Updates(User)
	return nil
}

// delete user
func DeleteUser(User *models.User, id int32) (err error) {
	database.Db.Where("user_id = ?", id).Delete(User)
	return nil
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
