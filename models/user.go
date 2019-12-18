package models

import "github.com/jinzhu/gorm"

type AuthSwag struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Model
	Username string `json:"username"`
	Password string `json:"password"`
	//Role     []Role `json:"role" gorm:"many2many:user_role;"`
}

func CheckUser(username, password string)  (bool,error) {
	var user User
	err := db.Select("id").Where(User{Username: username, Password: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true ,nil
	}

	return false ,nil
}