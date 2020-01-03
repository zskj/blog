package models

import "github.com/jinzhu/gorm"

//登录
type AuthSwag struct {
	Username    string `json:"username"` //登录账户
	Password    string `json:"password"` //登录密码
	CaptchaCode string `json:"captcha_code"`
	CaptchaId   string `json:"captcha_id"`
}

//注册
type Reg struct {
	Username      string `json:"username" binding:"required"`        //用户名
	Password      string `json:"password"  binding:"required"`       //密码
	PasswordAgain string `json:"password_again" binding:"required" ` //确认密码
	CaptchaCode   string `json:"captcha_code" binding:"required"`    //验证码
	CaptchaId     string `json:"captcha_id"  binding:"required"`     //验证码Id
}

// user 表
type User struct {
	Model
	Username  string `json:"username"`
	Password  string `json:"password"`
	Secret    string `json:"secret"`
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
//通过ID 查找用户
func FindUserById(id int) (User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}
//通过username 查找用户
func FindUserByUsername(username string) (User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, err
}
//创建新用户
func NewUser(user *User) (int, error) {
	err := db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}
