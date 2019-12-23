package user_service

import (
	"blog/models"
	"blog/pkg/util"
)

type User struct {
	ID       int
	Username string
	Password string
	Status   int
	CreatedOn  int
	ModifiedOn int
	DeletedOn  int

	PageNum  int
	PageSize int

}

func (a *User) Check() (bool, error) {
	return models.CheckUser(a.Username, util.EncodeMD5(a.Password))
}
