package models

import (
	"blog/pkg/util/hash"
	"blog/pkg/util/rand"
	"github.com/jinzhu/gorm"
)

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

//更新用户的secret
func UpdateUserSecret(user *User) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

//更新用户的secret
func UpdateUserNewPassword(user *User, newPassword string) (int, error) {
	var secretString string
	for {
		secretString = rand.RandStringBytesMaskImprSrcUnsafe(5)
		if user.Secret != secretString {
			break
		}
	}
	db.First(user)
	user.Secret = secretString
	user.Password = hash.EncodeMD5(newPassword)
	err := db.Save(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}
