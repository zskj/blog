package api

import (
	"blog/service/user_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"blog/pkg/app"
	"blog/pkg/e"
	"blog/pkg/util"
)

type auth struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role_id"`
}

// @Summary   获取登录token 信息
// @Tags 登录管理
// @Accept json
// @Produce  json
// @Param   body  body   models.AuthSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router /auth  [POST]
func Auth(c *gin.Context) {

	appG := app.Gin{C:c}
	var reqInfo auth
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.MaxSize(reqInfo.Username, 100, "username").Message("最长为100字符")
	valid.MaxSize(reqInfo.Password, 100, "password").Message("最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	    authService := user_service.User{Username: reqInfo.Username, Password: reqInfo.Password}
		isExist , err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(reqInfo.Username, reqInfo.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}