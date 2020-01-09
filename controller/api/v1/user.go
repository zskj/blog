package v1

import (
	"blog/models"
	"blog/models/schema"
	"blog/pkg/util/hash"
	"blog/pkg/util/rand"
	"github.com/dchest/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"blog/pkg/app"
	"blog/pkg/e"
	"blog/pkg/util"
)

type auth struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaCode string `json:"captcha_code"`
	CaptchaId   string `json:"captcha_id"`
}

type currentUser struct {
	Id string `json:"id"`

	Username string `json:"username"`
}

// @Summary   注册用户
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.Reg   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/reg  [POST]
func Reg(c *gin.Context) {

	appG := app.Gin{C: c}
	var reqInfo schema.Reg
	var data interface{}
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	valid := validation.Validation{}
	valid.Required(reqInfo.Username, "username").Message("请输入用户名")
	valid.Required(reqInfo.Password, "password").Message("输入密码")
	valid.Required(reqInfo.PasswordAgain, "password_again").Message("输入去确认密码")
	valid.MaxSize(reqInfo.Username, 50, "username").Message("密码最长为50字符")
	valid.MaxSize(reqInfo.Password, 50, "password").Message("密码最长为50字符")
	valid.MaxSize(reqInfo.PasswordAgain, 50, "password_again").Message("确认密码最长为50字符")
	valid.Required(reqInfo.CaptchaCode, "captcha_code").Message("请您输入验证码")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}
	if !captcha.VerifyString(reqInfo.CaptchaId, reqInfo.CaptchaCode) {
		appG.Response(http.StatusInternalServerError, e.ERROR_CAPTCHA_FAIL, data)
		return
	}
	if _, isExist := models.FindUserByUsername(reqInfo.Username); isExist != gorm.ErrRecordNotFound {
		appG.Response(http.StatusInternalServerError, e.ERROR_REPEAT_NAME, data)
		return
	}
	var newUser models.User
	newUser.Username = reqInfo.Username
	newUser.Password = hash.EncodeMD5(reqInfo.Password)
	newUser.Status = 1
	newUser.Secret = rand.RandStringBytesMaskImprSrcUnsafe(5)
	newUser.CreatedOn = int(time.Now().Unix())
	newUser.ModifiedOn = int(time.Now().Unix())
	userId, isSuccess := models.NewUser(&newUser)
	if userId > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"id": userId})
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, isSuccess)

}

// @Summary   用户登录 获取token 信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.AuthSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router /api/v1/auth  [POST]
func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo auth
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(reqInfo.Username, "username").Message("请输入用户名")
	valid.Required(reqInfo.Password, "password").Message("输入密码")
	valid.MaxSize(reqInfo.Username, 50, "username").Message("最长为50字符")
	valid.MaxSize(reqInfo.Password, 50, "password").Message("用户名最长为50字符")
	valid.Required(reqInfo.CaptchaCode, "captcha_code").Message("请您输入验证码")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}
	if !captcha.VerifyString(reqInfo.CaptchaId, reqInfo.CaptchaCode) {
		appG.Response(http.StatusInternalServerError, e.ERROR_CAPTCHA_FAIL, valid.Errors)
		return
	}

	isExist, user, err := models.LoginCheck(reqInfo.Username, hash.EncodeMD5(reqInfo.Password))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

// @Summary 获取登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/currentuser   [GET]
func CurrentUser(c *gin.Context) {
	var code int
	var data interface{}
	var user currentUser
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	//token := c.Query("token")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		user.Id = claims.Id
		user.Username = claims.Audience
	}

	if code != e.SUCCESS {
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"data": user,
		})
	}

}

// @Summary 刷新token
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/refreshtoken  [GET]
func RefreshToken(c *gin.Context) {
	var data interface{}
	var code int
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	if Authorization == "" {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	}
	token, err := util.RefreshToken(Authorization)
	if err != nil {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": err,
		})
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}

// @Summary 用户登出
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/logout  [POST]
func Logout(c *gin.Context) {
	var data interface{}
	var code int
	appG := app.Gin{C: c}
	code = e.SUCCESS
	claims := c.MustGet("claims").(*util.Claims)
	if claims == nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST, err)
		return
	}
	user, err := models.FindUserById(id)
	if err != nil {
		code = e.ERROR_EXIST_FAIL
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": err,
		})
	}
	_, isSuccess := models.UpdateUserSecret(&user)
	if isSuccess != nil {
		code = e.ERROR_EDIT_FAIL
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": isSuccess,
		})
	}
	appG.Response(http.StatusOK, code, map[string]interface{}{
		"data": data,
	})

}

// @Summary 登录用户修改密码
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   schema.PasswordSwag   true "body"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/password   [POST]
func Password(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo schema.PasswordSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	valid := validation.Validation{}
	valid.Required(reqInfo.NewPassword, "new_password").Message("请输入新密码")
	valid.Required(reqInfo.OldPassword, "old_password").Message("输入旧密码")
	valid.MaxSize(reqInfo.NewPassword, 50, "new_password").Message("新密码最长为50字符")
	valid.MaxSize(reqInfo.OldPassword, 50, "old_password").Message("密码最长为50字符")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}
	claims := c.MustGet("claims").(*util.Claims)
	if claims == nil {
		appG.Response(http.StatusOK, e.ERROR_AUTH, nil)
		return
	}
	id, err := strconv.Atoi(claims.Id)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST, err)
		return
	}
	user, err := models.FindUserById(id)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST, err)
		return
	}
	if hash.EncodeMD5(reqInfo.OldPassword) != user.Password {
		appG.Response(http.StatusBadRequest, e.INVALID_OLD_PASS, nil)
		return
	}
	_, isOk := models.UpdateUserNewPassword(&user, reqInfo.NewPassword)
	if isOk != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_EDIT_FAIL, isOk)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, isOk)
	return
}
