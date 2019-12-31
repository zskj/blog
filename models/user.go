package models

import "github.com/jinzhu/gorm"

type AuthSwag struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaCode string `json:"captcha_code"`
	CaptchaId   string `json:"captcha_id"`
}

type Reg struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	PasswordAgain string `json:"password_again"`
	CaptchaCode   string `json:"captcha_code"`
	CaptchaId     string `json:"captcha_id"`
}

type User struct {
	Model
	Username  string `json:"username"`
	Password  string `json:"password"`
	Status    int    `json:"status"`
	DeletedOn int    `json:"deleted_on"`
}

//登录验证
func LoginCheck(username, password string) (bool, User, error) {
	var user User
	err := db.Where(&User{Username: username, Password: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, user, err
	}
	if user.ID > 0 && user.Status > 0 {
		return true, user, nil
	}

	return false, user, nil
}

func FindUserById(id int) (User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}

func NewUser(user *User) (int, error) {
	err := db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}
